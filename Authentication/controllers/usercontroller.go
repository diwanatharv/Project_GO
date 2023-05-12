package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"awesomeProject2/Authentication/database"
	"awesomeProject2/Authentication/helpers"
	"awesomeProject2/Authentication/models"
)

var usercollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var Validator = validator.New()

func HashPassword(password string) string {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashpassword)
}
func VerifyPassword(userpassword string, providedpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userpassword), []byte(providedpassword))

}
func SingnUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationerr := Validator.Struct(user)
		if validationerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationerr.Error()})
			return
		}
		count, err := usercollection.CountDocuments(ctx, bson.M{"Email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking documents"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		count, err = usercollection.CountDocuments(ctx, bson.M{"Phone": user.Phone})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking documents"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()
		user.UserId = user.Id.Hex()
		token, refreshtoken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserId)
		user.Token = token
		user.RefreshToken = refreshtoken
		resNum, err := usercollection.InsertOne(ctx, user)
		if err != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resNum)
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := usercollection.FindOne(ctx, bson.M{"Email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}
		passwordisvalid, msg := VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()
		if passwordisvalid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if founduser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		token, refreshtoken, _ := helpers.GenerateAllTokens(*founduser.Email, *founduser.FirstName, *founduser.LastName, *founduser.UserType, founduser.UserId)
		helpers.UpdateAllTokens(token, refreshtoken, founduser.UserId)
		err = usercollection.FindOne(ctx, bson.M{"UserId": founduser.UserId}).Decode(&founduser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, founduser)
	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helpers.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordperpage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordperpage < 1 {
			recordperpage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		startindex := (page - 1) * recordperpage
		startindex, err = strconv.Atoi(c.Query("startIndex"))
		matchStage := bson.D{{"$match", bson.D{}}}
		groupstage := bson.D{{"$group", bson.D{{"Id", "null"}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectstage := bson.D{
			{"$project", bson.D{
				{"Id", 0},
				{"total_count", 1},
				{"user_items",bson.D{{"$slice",[]interface{}{"$data",startindex,recordperpage}}}}},
			}}
		}

res,err:=usercollection.Aggregate(ctx,mongo.Pipeline{
	matchStage,groupstage,projectstage
})

	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId := c.Param("UserId")
		err := helpers.MatchUserTypeToUid(c, UserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err = usercollection.FindOne(ctx, bson.M{UserId: user.UserId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}
		c.JSON(http.StatusOK, user)
	}
}
