package model

import (
	. "github.com/lucas-kern/tower-of-babel_server/app/model" // Import the package you want to test

	"testing"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"
	"os"
	"strings"
)

func TestNewBase(t *testing.T) {
    // Setup
    os.Setenv("BASE_LEVELS", "../../app/config/base_levels.json")

    // Teardown
    t.Cleanup(func() {
        os.Unsetenv("BASE_LEVELS")
    })

	owner := primitive.NewObjectID()
	base := NewBase(owner)

	assert.NotNil(t, base, "Expected base to be not nil")
	assert.False(t, base.ID.IsZero(), "Expected base ID to be not zero")
	assert.Equal(t, owner, base.Owner, "Expected base owner to be equal to input owner")
	assert.NotNil(t, base.Grid, "Expected base grid to be not nil")
	assert.Equal(t, 100, len(base.Grid), "Expected base grid to have 100 rows")
	assert.Equal(t, 100, len(base.Grid[0]), "Expected base grid to have 100 columns")
	assert.NotNil(t, base.PlacedBuildings, "Expected PlacedBuildings to be not nil")
	assert.NotNil(t, base.PendingBuildings, "Expected PendingBuildings to be not nil")

	expectedPendingBuildings := map[string][]Building{
			"tower": {
					Building{ Name: "tower" },
			},
			"sawmill": {
					Building{ Name: "sawMill" },
			},
			"woodstorage": {
				Building{ Name: "woodStorage" },
		},
	}

	assert.Equal(t, len(expectedPendingBuildings), len(base.PendingBuildings), "Expected pending buildings count to match")

	for key, buildings := range base.PendingBuildings {
			expectedBuildings, ok := expectedPendingBuildings[key]
			assert.True(t, ok, "Unexpected key in PendingBuildings: %s", key)

			assert.Equal(t, len(expectedBuildings), len(buildings), "Expected buildings count to match for key %s", key)

			for i, building := range buildings {
					assert.Equal(t, expectedBuildings[i].Name, building.Name, "Expected building at index %d for key %s to match", i, key)
			}
	}
}

func TestValidateBuildingPlacement(t *testing.T) {
	// Create a base with a grid and some pending buildings
	base := &Base{
			Grid: make([][]*Building, 10),
			PendingBuildings: map[string][]Building{
					"tower": {Building{Name: "tower"}},
			},
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, 10)
	}

	// Test a valid building placement
	building := &Building{
			Name: "tower",
			PosX: 5,
			PosZ: 5,
			Width: 2,
			Height: 2,
	}
	err := base.ValidateBuildingPlacement(building)
	assert.Nil(t, err, "Expected no error for valid building placement")

	// Test a building placement that is out of grid bounds
	building.PosX = 9
	building.PosZ = 9
	err = base.ValidateBuildingPlacement(building)
	assert.NotNil(t, err, "Expected error for building placement out of grid bounds")

	// Test a building placement that overlaps with an existing building
	base.Grid[5][5] = &Building{}
	building.PosX = 5
	building.PosZ = 5
	err = base.ValidateBuildingPlacement(building)
	assert.NotNil(t, err, "Expected error for building placement that overlaps with existing building")

	// Test a building placement where all buildings of this type have been placed
	base.PendingBuildings["tower"] = nil
	err = base.ValidateBuildingPlacement(building)
	assert.NotNil(t, err, "Expected error for building placement where all buildings of this type have been placed")
}

func TestValidateBuildingRemoval(t *testing.T) {
	// Create a base with a grid and a building
	base := &Base{
			Grid: make([][]*Building, 10),
			PlacedBuildings: make(map[string][]Building),
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, 10)
	}
	building := &Building{
			Name: "tower",
			PosX: 5,
			PosZ: 5,
			Width: 2,
			Height: 2,
	}
	base.Grid[5][5] = building
	base.Grid[5][6] = building
	base.Grid[6][5] = building
	base.Grid[6][6] = building
	base.PlacedBuildings[strings.ToLower(building.Name)] = append(base.PlacedBuildings[strings.ToLower(building.Name)], *building)

	// Test a valid building removal
	err := base.ValidateBuildingRemoval(building)
	assert.Nil(t, err, "Expected no error for valid building removal")

	// Test a building removal that is out of grid bounds
	building.PosX = 9
	building.PosZ = 9
	err = base.ValidateBuildingRemoval(building)
	assert.NotNil(t, err, "Expected error for building removal out of grid bounds")
	assert.Equal(t, fmt.Errorf("Building Placement failed: building is out of grid bounds").Error(), err.Error())

	// Test a building removal where the building is not at the given location
	building.PosX = 7
	building.PosZ = 7
	err = base.ValidateBuildingRemoval(building)
	assert.NotNil(t, err, "Expected error for building removal where building is not at given location")
	assert.Equal(t, fmt.Errorf("Building removal failed: Not same building in location").Error(), err.Error())

	// Test a building removal where no buildings of this type have been placed
	base.PlacedBuildings = map[string][]Building{}
	err = base.ValidateBuildingRemoval(building)
	assert.NotNil(t, err, "Expected error for building removal where no buildings of this type have been placed")
	assert.Equal(t, fmt.Errorf("No buildings of this type have been placed.").Error(), err.Error())
}

func TestAddBuildingToBase(t *testing.T) {
	// Create a base with a grid and some pending buildings
	base := &Base{
			Grid: make([][]*Building, 10),
			PendingBuildings: map[string][]Building{
					"tower": {Building{Name: "tower"}},
			},
			PlacedBuildings: make(map[string][]Building),
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, 10)
	}

	// Test adding a valid building
	building := &Building{
			Name: "tower",
			PosX: 5,
			PosZ: 5,
			Width: 2,
			Height: 2,
	}
	err := base.AddBuildingToBase(building)
	assert.Nil(t, err, "Expected no error for valid building addition")
	assert.NotNil(t, base.Grid[5][5], "Expected building to be placed on grid")
	assert.NotNil(t, base.Grid[5][6], "Expected building to be placed on grid")
	assert.NotNil(t, base.Grid[6][5], "Expected building to be placed on grid")
	assert.NotNil(t, base.Grid[6][6], "Expected building to be placed on grid")
	assert.NotNil(t, base.PlacedBuildings["tower"], "Expected building to be added to PlacedBuildings")
	assert.Equal(t, 0, len(base.PendingBuildings["tower"]), "Expected building to be removed from PendingBuildings")
	
	// Test adding a building that is out of grid bounds
	building.PosX = 9
	building.PosZ = 9
	err = base.AddBuildingToBase(building)
	assert.NotNil(t, err, "Expected error for building addition out of grid bounds")

	// Test adding a building where all buildings of this type have been placed
	base.PendingBuildings["tower"] = nil
	err = base.AddBuildingToBase(building)
	assert.NotNil(t, err, "Expected error for building addition where all buildings of this type have been placed")
}

func TestRemoveBuildingFromBase(t *testing.T) {
	// Create a base with a grid and a building
	base := &Base{
			Grid: make([][]*Building, 10),
			PlacedBuildings: make(map[string][]Building),
			PendingBuildings: make(map[string][]Building),
	}
	for i := range base.Grid {
			base.Grid[i] = make([]*Building, 10)
	}
	building := &Building{
			Name: "tower",
			PosX: 5,
			PosZ: 5,
			Width: 2,
			Height: 2,
	}
	base.Grid[5][5] = building
	base.Grid[5][6] = building
	base.Grid[6][5] = building
	base.Grid[6][6] = building
	base.PlacedBuildings[strings.ToLower(building.Name)] = append(base.PlacedBuildings[strings.ToLower(building.Name)], *building)

	// Test removing a valid building
	err := base.RemoveBuildingFromBase(building)
	assert.Nil(t, err, "Expected no error for valid building removal")
	assert.Nil(t, base.Grid[5][5], "Expected building to be removed from grid")
	assert.Nil(t, base.Grid[5][6], "Expected building to be removed from grid")
	assert.Nil(t, base.Grid[6][5], "Expected building to be removed from grid")
	assert.Nil(t, base.Grid[6][6], "Expected building to be removed from grid")
	assert.Equal(t, 0, len(base.PlacedBuildings["tower"]), "Expected building to be removed from PlacedBuildings")
	assert.Equal(t, 1, len(base.PendingBuildings["tower"]), "Expected building to be added to PendingBuildings")

	// Test removing a building that is out of grid bounds
	building.PosX = 9
	building.PosZ = 9
	err = base.RemoveBuildingFromBase(building)
	assert.NotNil(t, err, "Expected error for building removal out of grid bounds")

	// Test removing a building where no buildings of this type have been placed
	base.PlacedBuildings = map[string][]Building{}
	err = base.RemoveBuildingFromBase(building)
	assert.NotNil(t, err, "Expected error for building removal where no buildings of this type have been placed")
}


func TestGenerateNextLevelBuildings(t *testing.T) {

	// Setup
	os.Setenv("BASE_LEVELS", "../../app/config/base_levels.json")

	// Teardown
	t.Cleanup(func() {
			os.Unsetenv("BASE_LEVELS")
	})

	// Test generating buildings for a level with no available buildings
	buildings, err := GenerateNextLevelBuildings(1000) // Assuming level 1000 has no buildings
	assert.NotNil(t, err, "Expected error for level with no available buildings")
	assert.Nil(t, buildings, "Expected no buildings to be generated")

	// Test generating buildings for a negative level
	buildings, err = GenerateNextLevelBuildings(-1)
	assert.NotNil(t, err, "Expected error for negative level")
	assert.Nil(t, buildings, "Expected no buildings to be generated")

	// Test generating buildings for a valid level
	buildings, err = GenerateNextLevelBuildings(1)
	assert.Nil(t, err, "Expected no error for valid level")
	assert.NotEmpty(t, buildings, "Expected buildings to be generated")

	expectedBuildings := []*Building{
    &Building{ Name: "tower" },
    &Building{ Name: "sawMill" },
    &Building{ Name: "woodStorage" },
	}

	assert.Equal(t, len(expectedBuildings), len(buildings), "Expected buildings count to match")

	for i, buildings := range buildings {
		expectedBuilding := expectedBuildings[i]
		assert.Equal(t, expectedBuilding, buildings, "Expected building at index %d for key %d to match", i, i)
	}
}