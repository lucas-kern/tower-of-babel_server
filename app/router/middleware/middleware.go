package middleware

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// middleware is used to intercept incoming HTTP calls and apply general functions upon them.
func middleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
					log.Printf("HTTP request sent to %s from %s", r.URL.Path, r.RemoteAddr)

					// call next registered handler
					n(w, r, ps)
	}
}