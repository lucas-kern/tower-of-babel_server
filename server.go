package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Tower of Babel!")
	},)

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}