package model

import (
	// "fmt"
	// "log"
	// "math/rand"
	// "time"

	bson "go.mongodb.org/mongo-driver/bson/primitive"
)

// Building represents a simple building with location
type Building struct {
	Name		string		`json:"name,omitempty" bson:"name",omitempty"`
	PosX		float64 	`json:"posX,omitempty" bson:"posX,omitempty"`
	PosY		float64 	`json:"posY,omitempty" bson:"posY,omitempty"`
	PosZ		float64 	`json:"posZ,omitempty" bson:"posZ,omitempty"`
}

// Base represents a base owned by [User]s
type Base struct {
	ID          	bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Owner       	bson.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
	Tower      		Building				`json:"sphere,omitempty" bson:"sphere,omitempty"`
	ArmyCamp     	Building				`json:"cube,omitempty" bson:"cube,omitempty"`
	Barracks      Building				`json:"cylinder,omitempty" bson:"cylinder,omitempty"`
}