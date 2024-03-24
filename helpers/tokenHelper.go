package helper

import (
	"log"
	"os"
	"time"

	"example.com/m/v2/database"
	jwt "github.com/dgrijalva/jwt-go"
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

func GenerateAllTokens(email string, firstName string, lastName string, userType string, userId string) (token string, refreshToken string, err string) {
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
	refreshClaims := &SignDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(SECRET)
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString(SECRET)

	if err != "" {
		log.Panic(err)
		return
	}
	return token, refreshToken, err

}
