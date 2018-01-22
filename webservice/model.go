package main

// Person Type (more like an object)
type Person struct {
	ID        int      `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Email     string   `json:"email,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// Address Type
type Address struct {
	ID      int    `json:"id,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zipcode string `json:"zipcode,omitempty"`
}

// People - group of persons
var people []Person

// GameDefinition - Model for the base game
type GameDefinition struct {
}
