package model

import (
	// "fmt"
	// "log"
	// "math/rand"
	// "time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TODO extend buildings for each building type (storage, mines, etc...) Then replace the current base struct with arrays of those types so we can just go through each array in order to place all the buildings of that type.

// Building represents a simple building with location
type Building struct {
	Name		string		`json:"name,omitempty" bson:"name,omitempty"`
	IsPlaced 	bool		`json:"isPlaced" bson:"isPlaced"`
	PosX		float64 	`json:"posX" bson:"posX"`
	PosY		float64 	`json:"posY" bson:"posY"`
	PosZ		float64 	`json:"posZ" bson:"posZ"`
	Width 	float64		`json:"width,omitempty" bson:"width,omitempty"`
	Height 	float64		`json:"height,omitempty" bson:"height,omitempty"`
}

// Base represents a base owned by [User]s
type Base struct {
	ID          primitive.ObjectID `bson:"_id"`
	Owner       primitive.ObjectID `json:"owner,omitempty" bson:"owner,omitempty"`
	Buildings   []Building         `json:"buildings,omitempty" bson:"buildings,omitempty"`
}

func NewBase(user_id primitive.ObjectID) *Base {
	tower := Building{
			Name:     "Tower",
			IsPlaced: true,
			PosX:     0,
			PosY:     0,
			PosZ:     0,
	}

	return &Base{
			ID:        primitive.NewObjectID(),
			Owner:     user_id,
			Buildings: []Building{tower}, // Place the tower building in the Buildings slice
	}
}
