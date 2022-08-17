package router

import (
	"github.com/lucas-kern/tower-of-babel_server/app/router/handlers"
	"github.com/julienschmidt/httprouter"
)
// TODO setup github.com/julienschmidt/httprouter
func GetRouter() *httprouter.Router{
	router := httprouter.New()
	router.GET("/", handlers.Index)
	router.GET("/bases/:id", handlers.Bases)
	return router
}
