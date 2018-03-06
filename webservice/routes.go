package webservice

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
	// Front End
	Route{
		"WelcomePage",
		"GET",
		"/",
		WelcomePage,
	},
	Route{
		"AuthCallback",
		"GET",
		"/oauth2callback",
		AuthCallback,
	},

	// Back End
	Route{
		"GetPlayers",
		"GET",
		"/players",
		ValidateHandler(GetPlayers),
	},
	Route{
		"GetPlayer",
		"GET",
		"/players/{id}",
		ValidateHandler(GetPlayer),
	},
	Route{
		"CreatePlayer",
		"POST",
		"/createPlayer",
		ValidateHandler(CreatePlayer),
	},
	Route{
		"DeletePlayer",
		"DELETE",
		"/deletePlayer/{id}",
		ValidateHandler(DeletePlayer),
	},
}
