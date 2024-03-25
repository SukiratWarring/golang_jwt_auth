package helper

import (
	"fmt"
	"log"
	"os"
	"time"

	"example.com/m/v2/database"
	jwt "github.com/golang-jwt/jwt"
)

type SignDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserType  string `json:"userType"`
	UserID    string `json:"userId"`
	jwt.StandardClaims
}

var userCollection = database.OpenColletion(database.Client, "users")
var SECRET = os.Getenv(("SECRET"))

func GenerateAllTokens(email string, firstName string, lastName string, userType string, userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		UserID:    userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 20).Unix(),
		},
	}
	fmt.Println("aa", claims)
	refreshClaims := &SignDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(),
		},
	}
	var refreshToken string
	var token string
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(SECRET))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString([]byte(SECRET))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err

}
