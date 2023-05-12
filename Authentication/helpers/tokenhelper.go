package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"awesomeProject2/Authentication/database"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	UserType  string
	jwt.StandardClaims
}

var Usercollection *mongo.Collection = database.OpenCollection(database.Client, "User")
var Secret_key = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string, usertype string, userid string) (string, string, error) {
	claims := SignedDetails{
		Email:          email,
		FirstName:      firstname,
		LastName:       lastname,
		UserType:       usertype,
		Uid:            userid,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()},
	}
	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(Secret_key))
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(Secret_key))
	if err != nil {
		log.Fatal(err)
	}
	return token, refreshtoken, err
}
func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"Token", signedtoken})
	updateObj = append(updateObj, bson.E{"RefreshToken", signedrefreshtoken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"UpdatedAt", updated_at})
	upsert := true
	filter := bson.M{"UserId": userid}
	opt := options.UpdateOptions{
		Upsert: upsert,
	}
	_, err := Usercollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}
