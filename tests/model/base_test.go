package model

import (
	. "github.com/lucas-kern/tower-of-babel_server/app/model" // Import the package you want to test

	"testing"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewBase(t *testing.T) {
	// Create a user ID (you can replace this with your own user ID logic)
	userID := primitive.NewObjectID()

	// Create a new base
	base := NewBase(userID)

	// Check if the base ID is set
	if base.ID.IsZero() {
		t.Error("Expected base ID to be set, but it's zero.")
	}

	// Check if the owner ID is set correctly
	if base.Owner != userID {
		t.Errorf("Expected base owner to be %s, but got %s", userID, base.Owner)
	}

	// Check if the grid is initialized correctly
	gridWidth := 1000
	gridHeight := 1000
	if len(base.Grid) != gridHeight {
		t.Errorf("Expected grid height to be %d, but got %d", gridHeight, len(base.Grid))
	}
	if len(base.Grid[0]) != gridWidth {
		t.Errorf("Expected grid width to be %d, but got %d", gridWidth, len(base.Grid[0]))
	}

	// Check if the tower is placed at the correct position
	middleX := gridWidth / 2
	middleY := gridHeight / 2
	expectedPosX := float64(middleX) - 1.0 // Adjusted for tower width
	expectedPosY := float64(middleY) - 1.0 // Adjusted for tower height

	// Ensure the tower exists in the grid at the expected positions
	tower := base.Grid[int(expectedPosY)][int(expectedPosX)]
	if tower == nil {
		t.Error("Expected tower in the grid, but it's nil.")
	}

	// Check if the tower's properties are set correctly
	if tower.Name != "tower" {
		t.Errorf("Expected tower name to be 'tower', but got '%s'", tower.Name)
	}
	if !tower.IsPlaced {
		t.Error("Expected tower to be placed, but it's not.")
	}
	if tower.PosX != expectedPosX {
		t.Errorf("Expected tower PosX to be %f, but got %f", expectedPosX, tower.PosX)
	}
	if tower.PosY != expectedPosY {
		t.Errorf("Expected tower PosY to be %f, but got %f", expectedPosY, tower.PosY)
	}
	if tower.Width != 2.0 {
		t.Errorf("Expected tower Width to be 2.0, but got %f", tower.Width)
	}
	if tower.Height != 2.0 {
		t.Errorf("Expected tower Height to be 2.0, but got %f", tower.Height)
	}
}

func TestValidateBuildingPlacement_ValidPlacement(t *testing.T) {
	// Create a base with a grid and a building placed in a valid position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
			Grid: make([][]*Building, gridHeight),
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, gridWidth)
	}
	building := &Building{
			Name:    "tower",
			IsPlaced: true,
			PosX:    0,
			PosY:    0,
			Width:   2,
			Height:  2,
	}

	err := base.ValidateBuildingPlacement(building)

	assert.NoError(t, err, "Expected valid placement, but got an error")
}

func TestValidateBuildingPlacement_OutOfBounds(t *testing.T) {
	// Create a base with a grid and a building placed out of grid bounds
	gridWidth := 5
	gridHeight := 5
	base := &Base{
			Grid: make([][]*Building, gridHeight),
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, gridWidth)
	}
	building := &Building{
			Name:    "tower",
			IsPlaced: true,
			PosX:    4, // Placed at the right edge of the grid
			PosY:    0,
			Width:   2,
			Height:  2,
	}

	err := base.ValidateBuildingPlacement(building)

	assert.Error(t, err, "Expected error for out-of-bounds placement")
	assert.Contains(t, err.Error(), "out of grid bounds", "Expected error message about out-of-bounds placement")
}

func TestValidateBuildingPlacement_OverlapWithExistingBuilding(t *testing.T) {
	// Create a base with a grid and a building already placed in the target position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
			Grid: make([][]*Building, gridHeight),
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, gridWidth)
	}
	// Place an existing building in the target position
	existingBuilding := &Building{
			Name:    "tower",
			IsPlaced: true,
			PosX:    0,
			PosY:    0,
			Width:   2,
			Height:  2,
	}
	base.Grid[0][0] = existingBuilding
	building := &Building{
			Name:    "tower",
			IsPlaced: true,
			PosX:    0,
			PosY:    0,
			Width:   2,
			Height:  2,
	}

	err := base.ValidateBuildingPlacement(building)

	assert.Error(t, err, "Expected error for overlapping placement")
	assert.Contains(t, err.Error(), "overlaps with existing building", "Expected error message about overlapping placement")
}

func TestAddBuildingToBase_ValidPlacement(t *testing.T) {
	// Create a base with a grid and a building placed in a valid position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	building := &Building{
		Name:    "tower",
		IsPlaced: false, // Building is not placed initially
		PosX:    0,
		PosY:    0,
		Width:   2,
		Height:  2,
	}

	err := base.AddBuildingToBase(building)

	assert.NoError(t, err, "Expected valid placement, but got an error")
	assert.True(t, building.IsPlaced, "Expected building to be placed, but it's not")
}

func TestAddBuildingToBase_OutOfBounds(t *testing.T) {
	// Create a base with a grid and a building placed out of grid bounds
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	building := &Building{
		Name:    "tower",
		IsPlaced: false, // Building is not placed initially
		PosX:    0, // Placed at the right edge of the grid
		PosY:    4,
		Width:   2,
		Height:  2,
	}

	err := base.AddBuildingToBase(building)

	assert.Error(t, err, "Expected error for out-of-bounds placement")
	assert.False(t, building.IsPlaced, "Expected building not to be placed")
}

func TestAddBuildingToBase_OverlapWithExistingBuilding(t *testing.T) {
	// Create a base with a grid and an existing building placed in the target position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	// Place an existing building in the target position
	existingBuilding := &Building{
		Name:    "tower",
		IsPlaced: true,
		PosX:    0,
		PosY:    0,
		Width:   2,
		Height:  2,
	}
	base.Grid[0][0] = existingBuilding
	building := &Building{
		Name:    "tower",
		IsPlaced: false, // Building is not placed initially
		PosX:    0,
		PosY:    0,
		Width:   2,
		Height:  2,
	}

	err := base.AddBuildingToBase(building)

	assert.Error(t, err, "Expected error for overlapping placement")
	assert.False(t, building.IsPlaced, "Expected building not to be placed")
}

func TestAddBuildingToBase_BuildingInGrid(t *testing.T) {
	// Create a base with a grid
	gridWidth := 5
	gridHeight := 5
	base := &Base{
			Grid:      make([][]*Building, gridHeight),
			Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, gridWidth)
	}

	// Create a building to add to the base
	building := &Building{
			Name:    "tower",
			IsPlaced: false, // Building is not placed initially
			PosX:    1,
			PosY:    1,
			Width:   2,
			Height:  2,
	}

	// Add the building to the base
	err := base.AddBuildingToBase(building)
	assert.NoError(t, err, "Expected no error when adding building to the base")

	// Check if the building is present in the correct cells of the grid
	startX := int(building.PosX)
	startY := int(building.PosY)
	endX := startX + int(building.Width)
	endY := startY + int(building.Height)

	for i := startY; i < endY; i++ {
			for j := startX; j < endX; j++ {
					assert.Equal(t, building, base.Grid[i][j], "Expected building to be in the grid cell")
			}
	}
}

func TestAddBuildingToBase_BuildingInBuildingsMap(t *testing.T) {
	// Create a base with an empty grid
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid:     make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}

	// Create a building to add to the base
	building := &Building{
		Name:     "tower",
		IsPlaced: false, // Building is not placed initially
		PosX:     1,
		PosY:     1,
		Width:    2,
		Height:   2,
	}

	err := base.AddBuildingToBase(building)

	// Assertions
	assert.NoError(t, err, "Expected no error when adding building to base")
	assert.True(t, building.IsPlaced, "Expected building to be placed")
	assert.Contains(t, base.Buildings["tower"], *building, "Expected building to be in the Buildings map")
}

func TestRemoveBuildingFromBase_ValidRemoval(t *testing.T) {
	// Create a base with a grid and a building placed in a valid position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	// Place a building in the grid
	building := &Building{
		Name:    "tower",
		IsPlaced: true,
		PosX:    0,
		PosY:    0,
		Width:   1,
		Height:  1,
	}
	base.Buildings["tower"] = []Building{*building}
	base.Grid[0][0] = building

	err := base.RemoveBuildingFromBase(building)

	// Assertions
	assert.NoError(t, err, "Expected no error when removing building from base")
	assert.False(t, building.IsPlaced, "Expected building to be marked as not placed")
	assert.NotContains(t, base.Buildings["tower"], *building, "Expected building to be removed from Buildings map")
}

func TestRemoveBuildingFromBase_BuildingNotInGrid(t *testing.T) {
	// Create a base with an empty grid
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	// Create a building that is not in the grid
	building := &Building{
		Name:     "tower",
		IsPlaced: true,
		PosX:     0,
		PosY:     0,
		Width:    2,
		Height:   2,
	}

	err := base.RemoveBuildingFromBase(building)

	// Assertions
	assert.Error(t, err, "Expected error when removing a building not in the grid")
	assert.True(t, building.IsPlaced, "Expected building to remain placed")
}

func TestRemoveBuildingFromBase_ErrorInRemoveFromBuildings(t *testing.T) {
	// Create a base with a grid and a building placed in a valid position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	// Place a building in the grid
	building := &Building{
		Name:    "tower",
		IsPlaced: true,
		PosX:    0,
		PosY:    0,
		Width:   2,
		Height:  2,
	}
	base.Grid[0][0] = building

	err := base.RemoveBuildingFromBase(building)

	// Assertions
	assert.Error(t, err, "Expected error when removing from buildings fails")
	assert.True(t, building.IsPlaced, "Expected building to remain placed")
}

func TestRemoveBuildingFromBase_BuildingSameAsOther(t *testing.T) {
	// Create a base with a grid and a building placed in a valid position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	// Place a building in the grid
	building := &Building{
		Name:    "tower",
		IsPlaced: true,
		PosX:    0,
		PosY:    0,
		Width:   1,
		Height:  1,
	}
	base.Buildings["tower"] = []Building{*building}
	base.Grid[0][0] = building

	building2 := &Building{
		Name:    "tower",
		IsPlaced: false,
		PosX:    0,
		PosY:    0,
		Width:   1,
		Height:  1,
	}

	base.Buildings["tower"] = []Building{*building}
	base.Grid[0][0] = building

	err := base.RemoveBuildingFromBase(building2)

	// Assertions
	assert.NoError(t, err, "Expected no error when removing building from base")
	assert.False(t, building.IsPlaced, "Expected building to not be placed")
	assert.False(t, building2.IsPlaced, "Expected building2 to not be placed")
}

func TestRemoveBuildingFromBase_BuildingDifferentFromOther(t *testing.T) {
	// Create a base with a grid and a building placed in a valid position
	gridWidth := 5
	gridHeight := 5
	base := &Base{
		Grid: make([][]*Building, gridHeight),
		Buildings: make(map[string][]Building),
	}
	for i := range base.Grid {
		base.Grid[i] = make([]*Building, gridWidth)
	}
	// Place a building in the grid
	building := &Building{
		Name:    "tower",
		IsPlaced: true,
		PosX:    0,
		PosY:    0,
		Width:   1,
		Height:  1,
	}
	base.Buildings["tower"] = []Building{*building}
	base.Grid[0][0] = building

	building2 := &Building{
		Name:    "test",
		IsPlaced: true,
		PosX:    1,
		PosY:    0,
		Width:   1,
		Height:  1,
	}

	base.Buildings["tower"] = []Building{*building}
	base.Grid[0][0] = building

	err := base.RemoveBuildingFromBase(building2)

	// Assertions
	assert.Error(t, err, "Expected an error of wrong building")
	assert.True(t, building.IsPlaced, "Expected building to remain placed")
	assert.True(t, building2.IsPlaced, "Expected building2 to remain placed")
}
