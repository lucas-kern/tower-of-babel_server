package client

import (
	"github.com/lucas-kern/tower-of-babel_server/app/model"
)

// Base represents a base owned by [User]s
type ClientBase struct {
	Buildings 	map[string][]model.Building `json:"buildings,omitempty" bson:"buildings,omitempty"`
}

// NewClientBase sets up a client appropriate [model.Base]
func NewClientBase(base *model.Base) *ClientBase {
	return &ClientBase{
		Buildings: base.Buildings,
	}
}