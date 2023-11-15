package model

import (
	"fmt"
	"strings"

	"github.com/lucas-kern/tower-of-babel_server/app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Manages the base structs and methods 

// Base represents a base owned by [User]s
type Base struct {
	ID          			primitive.ObjectID 		`json:"id,omitempty" bson:"_id"`
	Owner       			primitive.ObjectID 		`json:"owner,omitempty" bson:"owner,omitempty"`
	PlacedBuildings 	map[string][]Building `json:"placedBuildings" bson:"placedBuildings,omitempty"`
	PendingBuildings	map[string][]Building `json:"pendingBuildings" bson:"pendingBuildings,omitempty"`
	Grid        			[][]*Building       	`json:"grid,omitempty" bson:"grid,omitempty"`
}

func NewBase(user_id primitive.ObjectID) *Base {
	// Create the grid and initialize it with nil values
	gridWidth := 100 // Set the width of the grid (adjust as needed)
	gridHeight := 100 // Set the height of the grid (adjust as needed)

	grid := make([][]*Building, gridHeight)
	for i := range grid {
		grid[i] = make([]*Building, gridWidth)
	}

	base := &Base{
		ID:        primitive.NewObjectID(),
		Owner:     user_id,
		Grid: 		 grid,
		PlacedBuildings: make(map[string][]Building),
		PendingBuildings: make(map[string][]Building),
	}

	buildings, err := GenerateNextLevelBuildings(1)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	base.addMultipleBuildingsToPendingBuildings(buildings)

	return base
}

// Validate that the building is able to be placed
func (base *Base) ValidateBuildingPlacement(building *Building) error {

	// Check if the map already has an entry for the building type
	buildingType := strings.ToLower(building.Name)

	if len(base.PendingBuildings[buildingType]) == 0 {
			// If the building type doesn't exist in pending buildings then there are no more the user can place
		return fmt.Errorf("All buildings of this type have been placed.")
	}

	// Check if base.Grid is nil or empty (no rows)
	if base.Grid == nil || len(base.Grid) == 0  {
		return fmt.Errorf("Base grid is not instantiated correctly")
	}

	gridSizeX := len(base.Grid[0])
	gridSizeZ := len(base.Grid)

	startX := int(building.PosX)
	startZ := int(building.PosZ)
	endX := startX + int(building.Width)
	endZ := startZ + int(building.Height)

	// Ensure the building is within the size of the grid
	if startX < 0 || startZ < 0 || endX > gridSizeX || endZ > gridSizeZ {
			return fmt.Errorf("Building Placement failed: building is out of grid bounds")
	}

	// Ensure there is no other building in the placement
	for i := startZ; i < endZ; i++ {
			for j := startX; j < endX; j++ {
					if base.Grid[i][j] != nil {
							return fmt.Errorf("Building Placement failed: building overlaps with existing building")
					}
			}
	}

	return nil
}

// Validate that the building is able to be placed
func (base *Base) ValidateBuildingRemoval(building *Building) error {
	gridSizeX := len(base.Grid[0])
	gridSizeZ := len(base.Grid)

	startX := int(building.PosX)
	startZ := int(building.PosZ)
	endX := startX + int(building.Width)
	endZ := startZ + int(building.Height)

	buildingType := strings.ToLower(building.Name)
	if len(base.PlacedBuildings[buildingType]) == 0 {
		// If the building type doesn't exist in placed buildings then there is not a building to remove
		return fmt.Errorf("No buildings of this type have been placed.")
	}

	// Ensure the building is within the size of the grid
	if startX < 0 || startZ < 0 || endX > gridSizeX || endZ > gridSizeZ {
			return fmt.Errorf("Building Placement failed: building is out of grid bounds")
	}

	// Ensure this building is at the given location
	for i := startZ; i < endZ; i++ {
			for j := startX; j < endX; j++ {
					if base.Grid[i][j] == nil || !base.Grid[i][j].Equal(building) {
							return fmt.Errorf("Building removal failed: Not same building in location")
					}
			}
	}

	return nil
}

// Add a Building to the base
func (base *Base) AddBuildingToBase(building *Building) error {
	if err := base.ValidateBuildingPlacement(building); err != nil {
		return err
	}
	base.removeFromPendingBuildings(building)
	base.addToPlacedBuildings(building)

	err := base.placeBuildingOnGrid(building)
	if err != nil {
		base.removeFromPlacedBuildings(building)
		return err
	}

	return nil
}

// Remove a Building from the base
func (base *Base) RemoveBuildingFromBase(building *Building) error {
	if err := base.ValidateBuildingRemoval(building); err != nil {
		return err
	}

	rb, err := base.removeFromPlacedBuildings(building)
	if err != nil {
		return err
	}
	base.addBuildingToPendingBuildings(rb)

	base.removeBuildingFromGrid(building)
	if err != nil {
		base.addToPlacedBuildings(building)
		return err
	}

	return nil
}

// Add a building to the grid of a base
func (base *Base) placeBuildingOnGrid(building *Building) error {
	startX := int(building.PosX)
	startZ := int(building.PosZ)
	endX := startX + int(building.Width)
	endZ := startZ + int(building.Height)

	for i := startZ; i < endZ; i++ {
			for j := startX; j < endX; j++ {
					base.Grid[i][j] = building
			}
	}

	return nil
}

// Remove a building from the grid of a base
func (base *Base) removeBuildingFromGrid(buildingToRemove *Building) (*Building) {
	startX := int(buildingToRemove.PosX)
	startZ := int(buildingToRemove.PosZ)
	endX := startX + int(buildingToRemove.Width)
	endZ := startZ + int(buildingToRemove.Height)

	removedBuilding := new(Building)

	for i := startZ; i < endZ; i++ {
		for j := startX; j < endX; j++ {			
			// Check if the current cell contains the building to remove
				// Set the cell to an empty building
				removedBuilding = base.Grid[i][j]
				base.Grid[i][j] = nil
		}
	}

	return removedBuilding
}

func (base *Base) addToPlacedBuildings(newBuilding *Building) {
	base.PlacedBuildings = addToBuildingsMap(newBuilding, base.PlacedBuildings)
}

func (base *Base) addMultipleBuildingsToPendingBuildings(newBuildings []*Building) {
	for _, building := range newBuildings {
		base.PendingBuildings = addToBuildingsMap(building, base.PendingBuildings)
	}
}

func (base *Base) addBuildingToPendingBuildings(newBuilding *Building) {
		base.PendingBuildings = addToBuildingsMap(&Building{ Name: newBuilding.Name}, base.PendingBuildings)
}

// Add a new building to a buildings map
func addToBuildingsMap(newBuilding *Building, buildings map[string][]Building) map[string][]Building{

	if buildings == nil {
		buildings = make(map[string][]Building)
	}
	// Check if the map already has an entry for the building type
	buildingType := strings.ToLower(newBuilding.Name)
	
	if existingBuildings, ok := buildings[buildingType]; ok {
			// If the building type exists, append the new building to the existing slice
			buildings[buildingType] = append(existingBuildings, *newBuilding)
	} else {
			// If the building type doesn't exist, create a new entry in the map
			buildings[buildingType] = []Building{*newBuilding}
	}

	return buildings
}

// Remove a building from the PlacedBuildings map
func (base *Base) removeFromPlacedBuildings(buildingToRemove *Building) (*Building, error) {
	removedBuilding, pb, err := removeFromBuildingsMap(buildingToRemove, base.PlacedBuildings)
	if err != nil {
		return nil, err
	}
	base.PlacedBuildings = pb
	return removedBuilding, err
}

// Remove a building from the PendingBuildings map
func (base *Base) removeFromPendingBuildings(buildingToRemove *Building) (*Building, error) {
	b := &Building{ Name: buildingToRemove.Name, PosX: 0, PosY: 0, PosZ: 0}
	removedBuilding, pb, err := removeFromBuildingsMap(b, base.PendingBuildings)
	if err != nil {
		return nil, err
	}
	base.PendingBuildings = pb
	return removedBuilding, err
}

// Remove a building from a Buildings map
func removeFromBuildingsMap(buildingToRemove *Building, buildings map[string][]Building) (*Building, map[string][]Building, error) {
	// Check if the map already has an entry for the building type
	buildingType := strings.ToLower(buildingToRemove.Name)
	if existingPlacedBuildings, ok := buildings[buildingType]; ok {
		for i, existingBuilding := range existingPlacedBuildings {
			// Check if the existing building is the one to be removed
			if existingBuilding.Equal(buildingToRemove) {
				// Remove the building by slicing it out of the slice
				buildings[buildingType] = append(existingPlacedBuildings[:i], existingPlacedBuildings[i+1:]...)
				return &existingPlacedBuildings[i], buildings, nil
			}
		}
	}
	// Building type not found in the map, return an error
	return nil, nil, fmt.Errorf("Building type '%s' not found in the map", buildingType)
}


// GenerateNextLevelBuildings generates the buildings available for the next level.
func GenerateNextLevelBuildings(currentLevel int) ([]*Building, error) {
	levelBuildings, err := helpers.ReadLevelBuildings()
	if err != nil {
		return nil, err
	}

	for _, lb := range levelBuildings {
		if lb.Level == currentLevel {
			// Create Building objects based on lb.Buildings
			// Add these Building objects to a slice and return it
			var nextLevelBuildings []*Building
			for _, buildingName := range lb.Buildings {
				b := &Building{ Name: buildingName}
				nextLevelBuildings = append(nextLevelBuildings, b)
			}
			return nextLevelBuildings, nil
		}
	}
	return nil, fmt.Errorf("No available buildings for the next level")
}