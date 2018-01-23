package main

import (
	"log"
	"net/http"
)

// main function to boot up everything
func main() {
	players = append(players, Player{ID: "1", Firstname: "John", Lastname: "Doe", Email: "john@funding.com", Phone: "5555555555", Address: &Address{ID: 1, City: "City X", State: "State X", Zipcode: "11101"}, CreatedOn: "20180122", Active: true, SignedUp: true})
	players = append(players, Player{ID: "2", Firstname: "Koko", Lastname: "Doe", Email: "koko@funding.com", Phone: "5555551234", Address: &Address{ID: 2, City: "City Z", State: "State Y", Zipcode: "11101"}, CreatedOn: "20180122", Active: false, SignedUp: false})

	var router = NewRouter()

	log.Fatal(http.ListenAndServe(":8000", router))
}
