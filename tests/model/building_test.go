package model

import (
		"testing"
		"github.com/lucas-kern/tower-of-babel_server/app/model"
)

func TestBuildingEqual_SameBuilding(t *testing.T) {
    building1 := model.Building{
        Name:    "tower",
        IsPlaced: true,
        PosX:    0,
        PosY:    0,
        Width:   2,
        Height:  2,
    }
    building2 := model.Building{
        Name:    "tower",
        IsPlaced: true,
        PosX:    0,
        PosY:    0,
        Width:   2,
        Height:  2,
    }

    if !building1.Equal(&building2) {
        t.Error("Expected building1 and building2 to be equal, but they are not.")
    }
}

func TestBuildingEqual_DifferentBuilding(t *testing.T) {
    building1 := model.Building{
        Name:    "tower",
        IsPlaced: true,
        PosX:    0,
        PosY:    0,
        Width:   2,
        Height:  2,
    }
    building2 := model.Building{
        Name:    "house",
        IsPlaced: true,
        PosX:    0,
        PosY:    0,
        Width:   2,
        Height:  2,
    }

    if building1.Equal(&building2) {
        t.Error("Expected building1 and building2 to be different, but they are equal.")
    }
}

func TestBuildingEqual_DifferentPosition(t *testing.T) {
    building1 := model.Building{
        Name:    "tower",
        IsPlaced: true,
        PosX:    0,
        PosY:    0,
        Width:   2,
        Height:  2,
    }
    building2 := model.Building{
        Name:    "tower",
        IsPlaced: true,
        PosX:    1,
        PosY:    1,
        Width:   2,
        Height:  2,
    }

    if building1.Equal(&building2) {
        t.Error("Expected building1 and building2 to be different due to different positions, but they are equal.")
    }
}