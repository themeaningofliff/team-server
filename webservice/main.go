package main

import (
	"log"
	"net/http"
)

// main function to boot up everything
func main() {
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})

	var router = NewRouter()

	log.Fatal(http.ListenAndServe(":8000", router))
}
