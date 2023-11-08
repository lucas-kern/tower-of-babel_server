package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all account objects retrieved or inserted into the DB
//TODO add back the bson
type User struct {
    ID            primitive.ObjectID `bson:"_id"`
    FirstName    *string            `json:"firstName" validate:"required,min=2,max=100"`
    LastName     *string            `json:"lastName" validate:"required,min=2,max=100"`
    Password      *string           `json:"password" validate:"required,min=6"`
    Email         *string           `json:"email" validate:"email,required"`
    Token         *string           `json:"token"`
    RefreshToken *string            `json:"refreshToken"`
    CreatedAt    time.Time          `json:"createdAt"`
    UpdatedAt    time.Time          `json:"updatedAt"`
}

//ClientUser is the model that governs all account objects retrieved or inserted into the DB
type ClientUser struct {
	FirstName    *string            `json:"firstName,omitempty"`
	LastName     *string            `json:"lastName,omitempty"`
	Token         *string            `json:"token,omitempty"`
    RefreshToken *string            `json:"refreshToken,omitempty"`
    Base          *Base              `json:"Base,omitempty"`
}

// newUser sets up a client appropriate [model.User]
func NewUser(user *User) *ClientUser {
	return &ClientUser{
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Token:			 user.Token,
		RefreshToken:   user.RefreshToken,
	}
}