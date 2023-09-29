package model

import (
	"testing"
	"github.com/lucas-kern/tower-of-babel_server/app/model" // Import the package you want to test


	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewBase(t *testing.T) {
	// Create a user ID (you can replace this with your own user ID logic)
	userID := primitive.NewObjectID()

	// Create a new base
	base := model.NewBase(userID)

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

	tower = base.Grid[int(expectedPosY+1)][int(expectedPosX)]
	if tower == nil {
		t.Error("Expected tower in the grid, but it's nil.")
	}

	tower = base.Grid[int(expectedPosY)][int(expectedPosX+1)]
	if tower == nil {
		t.Error("Expected tower in the grid, but it's nil.")
	}

	tower = base.Grid[int(expectedPosY+1)][int(expectedPosX+1)]
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