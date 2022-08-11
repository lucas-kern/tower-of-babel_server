package router

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func HandleRoute(route string) func(w http.ResponseWriter, r *http.Request){

	//TODO merge with middleware to handle this better
	switch route {
	case "/":
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome to Tower of Babel!")
		}
	case "/bases":
		return func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(bases[0])
		}
	default:
		//TODO raise an error and handle correctly
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "That is not a correct call!")
		}
	}
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
    {ID: 1, Name: "Coffee's Base", Sphere: []int{4,3}, Cube: []int{24,25}, Cylinder: []int{8,9}},
    {ID: 1, Name: "Lucas' Base", Sphere: []int{4,3}, Cube: []int{24,25}, Cylinder: []int{8,9}},
}