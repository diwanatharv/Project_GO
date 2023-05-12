package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	usertype := c.GetString("UserType")
	err = nil
	if usertype != role {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
}
func MatchUserTypeToUid(c *gin.Context, UserId string) (err error) {
	usertype := c.GetString("UserType")
	uid := c.GetString("uid")
	err = nil
	if usertype == "USER" && uid != UserId {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, usertype)
	return err
}
