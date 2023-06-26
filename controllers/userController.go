package controllers

import (
	"context"
	"net/http"
	"time"

	"example.com/m/v2/database"
	helper "example.com/m/v2/helpers"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Client = database.OpenColletion(database.Client, "user")
var validate = validator.New()

func HashPassword() {

}
func VerifyPassword() {

}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

	}

}
func Login() {

}
func GetUsers() {

}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Params("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)

	}
}
