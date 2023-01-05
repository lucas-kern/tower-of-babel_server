package middleware

import (
		"fmt"
		"log"
    "net/http"

		"github.com/julienschmidt/httprouter"
    "github.com/lucas-kern/tower-of-babel_server/app/auth"
)

// middleware is used to intercept incoming HTTP calls and apply general functions upon them.
func middleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
					log.Printf("HTTP request sent to %s from %s", r.URL.Path, r.RemoteAddr)

					// call registered handler
					n(w, r, ps)
	}
}

// Authz validates token and authorizes users
// TODO update to use current router
func Authentication(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        clientToken := r.Header.Get("token")
        if clientToken == "" {
            // c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
						// c.Abort()
						log.Printf("There is no token")

            return
        }

        claims, err := auth.ValidateToken(clientToken)
        if err != "" {
						log.Panic(err)
            return
        }

				c := r.Context
        c.Set("email", claims.Email)
        c.Set("first_name", claims.First_name)
        c.Set("last_name", claims.Last_name)
        c.Set("uid", claims.Uid)

        r.WithContext(c)

				n(w, r, ps)
    }
}