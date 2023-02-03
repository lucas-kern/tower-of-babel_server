package handlers

import (
	"github.com/lucas-kern/tower-of-babel_server/app/server/database"
)

// HandlerEnv is a wrapper for the genral request handling and contains a database instance
type HandlerEnv struct {
	database *database.Database
}

// NewHandlerEnv returns a new [HandlerEnv] with the specified database
func NewHandlerEnv(db *database.Database) *HandlerEnv {
	return &HandlerEnv{
		database: db,
	}
}