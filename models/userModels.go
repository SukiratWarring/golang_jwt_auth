package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required" min="2" max="100""`
	Last_name     *string            `json:"last_name" validate:"required" min="2" max="100""`
	Token         *string            `json:"token"`
	Email         *string            `json:"email" vaildate:"required""`
	User_type     *string            `json:"user_type" validate :"required, eq=ADMIN || eq=USER"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_id"`
	Updated_at    time.Time          `json:updated_id`
	User_id       string             `json:user_id`
	Password      string             `json:password`
}
