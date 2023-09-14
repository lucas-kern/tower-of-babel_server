package router

import (
	"github.com/lucas-kern/tower-of-babel_server/app/router/handlers"
	"github.com/lucas-kern/tower-of-babel_server/app/router/middleware"
	"github.com/lucas-kern/tower-of-babel_server/app/server/database"
	"github.com/julienschmidt/httprouter"
)

// This registers all our routes and can wrap them in middle ware for auth and other items
// Returns the router with paths and handlers
//TODO add a refresh token endpoint and route
func GetRouter(db *database.Database) *httprouter.Router{
	EnvHandler := handlers.NewHandlerEnv(db)
	router := httprouter.New()
	router.GET("/", EnvHandler.Index)
	router.GET("/bases/:id", middleware.Authentication(EnvHandler.Bases))
	router.POST("/bases/:id", middleware.Authentication(EnvHandler.Bases))
	router.POST("/users/signup", EnvHandler.SignUp)
	router.POST("/users/login", EnvHandler.Login)
	router.POST("/token", EnvHandler.TokenRefresh)
	return router
}
