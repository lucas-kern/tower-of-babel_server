package router

import (
	"github.com/lucas-kern/tower-of-babel_server/app/router/handlers"
	"github.com/lucas-kern/tower-of-babel_server/app/router/middleware"
	"github.com/julienschmidt/httprouter"
)

// Returns the router with paths and handlers
//TODO pass the database to the signup handler
func GetRouter(db *database.Database) *httprouter.Router{
	router := httprouter.New()
	router.GET("/", handlers.Index)
	router.GET("/bases/:id", handlers.Bases)
	router.POST(/users/signup, handlers.SignUp)
	// router.POST(/users/login, middleware.Authentication(handlers.Login))
	return router
}
