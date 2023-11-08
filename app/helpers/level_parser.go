package helpers

import (
	"os"
	"encoding/json"
	"fmt"
)

// LevelBuilding represents the structure of available buildings for a specific level.
type LevelBuilding struct {
	Level     int      `json:"level"`
	Buildings []string `json:"buildings"`
}

// ReadLevelBuildings reads the JSON file and returns a slice of LevelBuilding objects.
func ReadLevelBuildings() ([]LevelBuilding, error) {
	filename := os.Getenv("BASE_LEVELS")
	if filename == "" {
		return nil, fmt.Errorf("Environment variable `BASE_LEVELS` not set")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var levelBuildings []LevelBuilding
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&levelBuildings); err != nil {
		return nil, err
	}

	return levelBuildings, nil
}