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
		json.NewEncoder(w).Encode(bases)
	})

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type base struct {
    ID     int  `json:"id"`
		Name   string  `json:"name"`
		Bases []string  `json:"bases"`
}

var bases = []base{
    {ID: 1, Name: "bilingual base", Bases: []string{"sphere", "rectangle", "square"}},
    {ID: 2, Name: "idiom base", Bases: []string{"triangle", "square", "square"}},
    {ID: 3, Name: "dialect base", Bases: []string{"square", "triangle", "rectangle"}},
}

