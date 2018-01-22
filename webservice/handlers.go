package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetPlayers displays all from the people var
func GetPlayers(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(players)
}

// GetPlayer displays a single data
func GetPlayer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	for _, item := range players {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Player{})
}

// CreatePlayer creates a new item
func CreatePlayer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var player Player
	_ = json.NewDecoder(r.Body).Decode(&player)
	player.ID = params["id"]
	players = append(players, player)
	json.NewEncoder(w).Encode(players)
}

// DeletePlayer deletes an item
func DeletePlayer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	for index, item := range players {
		if item.ID == params["id"] {
			players = append(players[:index], players[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(players)
	}
}
