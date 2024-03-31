package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, roleType string) (err error) {
	userType := c.Query("role")
	err = nil
	if userType != roleType {
		return errors.New("Unauthorized to access the resource")
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access the resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
