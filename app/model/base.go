package model

import (
	"fmt"
	"log"
	"strings"
	// "math/rand"
	// "time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TODO extend buildings for each building type (storage, mines, etc...) Then replace the current base struct with arrays of those types so we can just go through each array in order to place all the buildings of that type.

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

// Base represents a base owned by [User]s
type Base struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Owner       primitive.ObjectID `json:"owner,omitempty" bson:"owner,omitempty"`
	Buildings 	map[string][]Building `json:"buildings,omitempty" bson:"buildings,omitempty"`
	Grid        [][]*Building       `json:"grid,omitempty" bson:"grid,omitempty"`
}

func NewBase(user_id primitive.ObjectID) *Base {
	tower := Building{
		Name:     "tower",
		IsPlaced: true,
		PosX:     0,
		PosY:     0,
		PosZ:     0,
		Width:    2, // Set the width of the tower (adjust as needed)
		Height:   2, // Set the height of the tower (adjust as needed)
	}

	// Create the grid and initialize it with nil values
	gridWidth := 1000 // Set the width of the grid (adjust as needed)
	gridHeight := 1000 // Set the height of the grid (adjust as needed)

	grid := make([][]*Building, gridHeight)
	for i := range grid {
		grid[i] = make([]*Building, gridWidth)
	}

	// Calculate the middle of the grid
	middleX := gridWidth / 2
	middleY := gridHeight / 2

	// Calculate the initial position of the tower to place it in the middle
	tower.PosX = float64(middleX) - tower.Width/2
	tower.PosY = float64(middleY) - tower.Height/2

	base := &Base{
		ID:        primitive.NewObjectID(),
		Owner:     user_id,
	}

	// Add the tower to the building
	err := base.AddBuildingToBase(&tower)
	if err != nil {
		log.Panic(err)
	}

	// Place the tower in the grid
	err = base.placeBuildingOnGrid(&tower)
	if err != nil {
		log.Panic(err)
	}

	return base
}

// Validate that the building is able to be placed
func (base *Base) validateBuildingPlacement(building *Building) error {
	gridSizeX := len(base.Grid[0])
	gridSizeY := len(base.Grid)

	startX := int(building.PosX)
	startY := int(building.PosY)
	endX := startX + int(building.Width)
	endY := startY + int(building.Height)

	//TODO ensure that the amount of that type of building is not more than the user is allowed
	// Add a map of building types to confirm it is a valid building type then compare the level of the user's base to how many of each building they are allowed to use
	// buildingName := building.Name

	// barracksCount := len(base.Buildings[buildingName])

	// Ensure the building is within the size of the grid
	if startX < 0 || startY < 0 || endX > gridSizeX || endY > gridSizeY {
			return fmt.Errorf("Building Placement failed: building is out of grid bounds")
	}

	// Ensure there is no other building in the placement
	for i := startY; i < endY; i++ {
			for j := startX; j < endX; j++ {
					if base.Grid[i][j] != nil {
							return fmt.Errorf("Building Placement failed: building overlaps with existing building")
					}
			}
	}

	return nil
}

// Add a building to the grid of a base
func (base *Base) placeBuildingOnGrid(building *Building) error {
	startX := int(building.PosX)
	startY := int(building.PosY)
	endX := startX + int(building.Width)
	endY := startY + int(building.Height)

	for i := startY; i < endY; i++ {
			for j := startX; j < endX; j++ {
					base.Grid[i][j] = building
			}
	}

	return nil
}

// Remove a building from the grid of a base
func (base *Base) removeBuildingFromGrid(buildingToRemove *Building) error {
	startX := int(buildingToRemove.PosX)
	startY := int(buildingToRemove.PosY)
	endX := startX + int(buildingToRemove.Width)
	endY := startY + int(buildingToRemove.Height)

	emptyBuilding := new(Building)

	for i := startY; i < endY; i++ {
		for j := startX; j < endX; j++ {			
			// Check if the current cell contains the building to remove
			if base.Grid[i][j].Equal(buildingToRemove) {
				// Set the cell to an empty building
				base.Grid[i][j] = emptyBuilding
			}
		}
	}

	return nil
}

// Add a new building to the Buildings map
func (base *Base) addToBuildings(newBuilding *Building) error {
	// Check if the map already has an entry for the building type
	buildingType := strings.ToLower(newBuilding.Name)
	if existingBuildings, ok := base.Buildings[buildingType]; ok {
			// If the building type exists, append the new building to the existing slice
			base.Buildings[buildingType] = append(existingBuildings, *newBuilding)
	} else {
			// If the building type doesn't exist, create a new entry in the map
			base.Buildings[buildingType] = []Building{*newBuilding}
	}
	return nil
}

// Remove a building from the Buildings map
func (base *Base) removeFromBuildings(buildingToRemove *Building) error {
	// Check if the map already has an entry for the building type
	buildingType := strings.ToLower(buildingToRemove.Name)
	if existingBuildings, ok := base.Buildings[buildingType]; ok {
		for i, existingBuilding := range existingBuildings {
			// Check if the existing building is the one to be removed
			if existingBuilding.Equal(buildingToRemove) {
				// Remove the building by slicing it out of the slice
				base.Buildings[buildingType] = append(existingBuildings[:i], existingBuildings[i+1:]...)
				return nil
			}
		}
	}

	return nil
}

// Add a Building to the base
//TODO finish this and test it 
func (base *Base) AddBuildingToBase(building *Building) error {
	if err := base.validateBuildingPlacement(building); err != nil {
		return err
	}

	err := base.addToBuildings(building)
	if err != nil {
		return err
	}

	err = base.placeBuildingOnGrid(building)
	if err != nil {
		base.removeFromBuildings(building)
		return err
	}

	building.IsPlaced = true
	return nil
}

// Remove a Building from the base
//TODO finish this and test it 
func (base *Base) RemoveBuildingFromBase(building *Building) error {
	if err := base.validateBuildingPlacement(building); err != nil {
		return err
	}

	err := base.removeFromBuildings(building)
	if err != nil {
		return err
	}

	err = base.removeBuildingFromGrid(building)
	if err != nil {
		base.addToBuildings(building)
		return err
	}

	building.IsPlaced = true
	return nil
}

// TODO create method that adds building to buildings array and to grid. And a method that will remove it from each. Need to make a method that does everything that is needed when a building is added to the base