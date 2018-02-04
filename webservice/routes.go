package webservice

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

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
		"GetPeople",
		"GET",
		"/people",
		ValidateHandler(GetPeople),
	},
	Route{
		"GetPerson",
		"GET",
		"/people/{id}",
		ValidateHandler(GetPerson),
	},
	Route{
		"CreatePerson",
		"POST",
		"/createPerson",
		ValidateHandler(CreatePerson),
	},
	Route{
		"DeletePerson",
		"DELETE",
		"/people/{id}",
		ValidateHandler(DeletePerson),
	},
}
