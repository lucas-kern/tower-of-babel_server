package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"

	"github.com/julienschmidt/httprouter"
	. "github.com/lucas-kern/tower-of-babel_server/app/model"
)

//TODO merge with middleware to handle authentication
// TODO remove this Index
	func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Welcome to Tower of Babel!")
	}
	func Bases(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			id, err := strconv.Atoi(params.ByName("id"))
			if err != nil {return};
			json.NewEncoder(w).Encode(bases[id])
	}

var bases = []Base{
}
