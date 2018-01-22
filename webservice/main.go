package main

import (
	"log"
	"net/http"
)

// main function to boot up everything
func main() {
	players = append(players, Player{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	players = append(players, Player{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})

	var router = NewRouter()

	log.Fatal(http.ListenAndServe(":8000", router))
}
