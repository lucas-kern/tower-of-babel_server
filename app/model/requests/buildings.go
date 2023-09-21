package model

import (
	// "fmt"
	// "log"
	// "math/rand"
	// "time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// Building represents a simple building with location
type BuildingPlacement struct {
	Name			string		`json:"name,omitempty" bson:"name,omitempty"`
	PosX			float64 	`json:"posX" bson:"posX"`
	PosY			float64 	`json:"posY" bson:"posY"`
	PosZ			float64 	`json:"posZ" bson:"posZ"`
	Width 		float64		`json:"width,omitempty" bson:"width,omitempty"`
	Height 		float64		`json:"height,omitempty" bson:"height,omitempty"`
}