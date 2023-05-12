package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// eq concept is just like enum
type User struct {
	Id           primitive.ObjectID `json:"Id" bson:"Id"`
	FirstName    *string            `json:"FirstName" bson:"FirstName" validate:"required,min=2,max=50"`
	LastName     *string            `json:"LastName" bson:"LastName" validate:"required,min=2,max=50"`
	Password     *string            `json:"Password" bson:"Password" validate:"required,min=6"`
	Email        *string            `json:"Email" bson:"Email" validate:"email,required"`
	Phone        *string            `json:"Phone" bson:"Phone" validate:"required"`
	Token        *string            `json:"Token" bson:"Token"`
	UserType     *string            `json:"UserType" bson:"UserType" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken *string            `json:"RefreshToken" bson:"RefreshToken"`
	CreatedAt    time.Time          `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt    time.Time          `json:"UpdatedAt" bson:"UpdatedAt"`
	UserId       string             `json:"UserId" bson:"UserId"`
}
