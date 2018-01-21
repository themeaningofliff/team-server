package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetPeople displays all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(people)
}

// GetPerson displays a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// CreatePerson creates a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// DeletePerson deletes an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
