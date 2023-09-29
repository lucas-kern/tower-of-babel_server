package model

// Manages the building struct and methods 

// Building represents a simple building with location
type Building struct {
	Name			string		`json:"name,omitempty" bson:"name,omitempty"`
	IsPlaced 	bool			`json:"isPlaced" bson:"isPlaced"`
	PosX			float64 	`json:"posX" bson:"posX"`
	PosY			float64 	`json:"posY" bson:"posY"`
	PosZ			float64 	`json:"posZ" bson:"posZ"`
	Width 		float64		`json:"width,omitempty" bson:"width,omitempty"`
	Height 		float64		`json:"height,omitempty" bson:"height,omitempty"`
}

// Equal checks if two Building instances are equal.
func (b *Building) Equal(other *Building) bool {
	// Compare all relevant fields for equality
	return b.Name == other.Name && b.PosX == other.PosX && b.PosY == other.PosY && b.Width == other.Width && b.Height == other.Height
}