package main

import (
	"net/http"
)

// Route - defines struct of routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - a collection of Route
type Routes []Route

var routes = Routes{
	Route{
		"GetPlayers",
		"GET",
		"/players",
		GetPlayers,
	},
	Route{
		"GetPlayer",
		"GET",
		"/players/{id}",
		GetPlayer,
	},
	Route{
		"CreatePlayer",
		"POST",
		"/players/{id}",
		CreatePlayer,
	},
	Route{
		"DeletePlayer",
		"DELETE",
		"/people/{id}",
		DeletePlayer,
	},
}
