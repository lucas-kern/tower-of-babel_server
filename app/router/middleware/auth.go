package middleware

import (
		"encoding/json"
		"log"
    "net/http"

		"github.com/julienschmidt/httprouter"
		"github.com/lucas-kern/tower-of-babel_server/app/auth"
		"github.com/lucas-kern/tower-of-babel_server/app/model"
)

// A middleware that will take a token from the header and ensure this user is valid
// Authentication validates token and authorizes users
//TODO add a method for refreshing the token
func Authentication(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        clientToken := r.Header.Get("Authorization")
        if clientToken == "" {
						log.Printf("There is no authorization token")

            return
				}

        claims, err := auth.ValidateToken(clientToken)
        if err != "" {
						log.Panic(err)
            return
        }

				//TODO this works, but we don't need to send the user data back every time we authenticate
				var user model.User
				user.Email = &claims.Email
				user.First_name = &claims.First_name
				user.Last_name = &claims.Last_name
				user.User_id = claims.Uid

				json.NewEncoder(w).Encode(user)

				n(w, r, ps)
    }
}