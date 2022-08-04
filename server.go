package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Tower of Babel!")
	},)

	http.HandleFunc("/bases", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(bases[0])
	})

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
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

