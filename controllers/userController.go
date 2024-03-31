package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/m/v2/database"
	helper "example.com/m/v2/helpers"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenColletion(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(byte)

}
func VerifyPassword(userPass string, providePass string) (bool, string) {
	isValid := true
	msg := ""
	err := bcrypt.CompareHashAndPassword([]byte(providePass), []byte(userPass))
	if err != nil {
		isValid = false
		fmt.Println(err)
		msg = fmt.Sprintf("Email or password is incorrect.")
	}
	return isValid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		// parse the JSON payload of the request body and bind it to a Go struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("user", user.ID.Hex())
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		fmt.Println("Signing up: ")
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		fmt.Println(count)
		password := HashPassword(user.Password)
		user.Password = password
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		// user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		// user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		instertionNumber, err := userCollection.InsertOne(ctx, user)
		if (err) != nil {
			msg := fmt.Sprintf("Error while inserting user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, instertionNumber)

	}

}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// creating a new context with a timeout using the context package in Go.
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		if err != nil {
			msg := fmt.Sprintf("Error while finding user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		fmt.Println(user.Password, founduser.Password)
		value, msg := VerifyPassword(user.Password, founduser.Password)
		if value == false {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		token, refreshToken, err := helper.GenerateAllTokens(*founduser.Email, *founduser.First_name, *founduser.Last_name, *founduser.User_type, *&founduser.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating token"})
			return
		}
		defer cancel()
		//Update all tokens
		opts := options.FindOneAndUpdate().SetUpsert(true)
		err = userCollection.FindOneAndUpdate(
			ctx,
			bson.M{"email": user.Email}, // Filter document
			bson.M{"$set": bson.M{"token": token, "refreshToken": refreshToken}}, // Update document
			opts,
		).Decode(&founduser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while Updating token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":        token,
			"refreshToken": refreshToken,
		})
		return

	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helper.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		recordPerPage, err := strconv.Atoi(c.Query("pageSize"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		pageNo, err := strconv.Atoi(c.Query("pageNo"))
		if err != nil || pageNo < 1 {
			pageNo = 1
		}

		// creating a new context with a timeout using the context package in Go.
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		opts := options.Find().SetSort(bson.D{{"updated_id", -1}}).SetLimit(int64(recordPerPage)).SetSkip((int64(pageNo - 1)))
		cursor, err := userCollection.Find(ctx, bson.M{}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Get a list of all returned documents and print them out.
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{"All Users": results})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

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
