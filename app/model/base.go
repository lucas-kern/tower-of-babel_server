package model

import (
	// "fmt"
	// "log"
	// "math/rand"
	// "time"

	"gopkg.in/mgo.v2/bson"
)

// Building represents a simple building with location
type Building struct {
	PosX				float64 		`json:"posX,omitempty" bson:"posX,omitempty"`
	PosY				float64 		`json:"posY,omitempty" bson:"posY,omitempty"`
	PosZ				float64 		`json:"posZ,omitempty" bson:"posZ,omitempty"`
}

// Base represents a base owned by [User]s
type Base struct {
	ID          	bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty"`
	Owner       	bson.ObjectId   `json:"owner,omitempty" bson:"owner,omitempty"`
	Tower      		Building				`json:"sphere,omitempty" bson:"sphere,omitempty"`
	ArmyCamp     	Building				`json:"cube,omitempty" bson:"cube,omitempty"`
	Barracks      Building				`json:"cylinder,omitempty" bson:"cylinder,omitempty"`
}