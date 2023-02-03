
package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// TODO remove this Index
	func (env *HandlerEnv) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Welcome to Tower of Babel!")
	}
	func (env *HandlerEnv) Bases(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			id, err := strconv.Atoi(params.ByName("id"))
			if err != nil {return};
			json.NewEncoder(w).Encode(bases[id])
	}

type base struct {
    ID     int  `json:"id"`
		Name   string  `json:"name"`
		Sphere []int  `json:"sphere"`
		Cube []int  `json:"cube"`
		Cylinder []int  `json:"cylinder"`
}

var bases = []base{
    {ID: 1, Name: "Diego's Base", Sphere: []int{4,3}, Cube: []int{2,3}, Cylinder: []int{0,1}},
    {ID: 2, Name: "Coffee's Base", Sphere: []int{4,3}, Cube: []int{24,25}, Cylinder: []int{8,9}},
    {ID: 3, Name: "Lucas' Base", Sphere: []int{4,3}, Cube: []int{24,25}, Cylinder: []int{8,9}},
}