package webservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	oauth2 "golang.org/x/oauth2"
	// Google API Client
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

var state string
var store = sessions.NewCookieStore([]byte("secret")) // TODO: This probably needs to be changed.

// WelcomePage greets a user who has found the site and presents them with a Login with Google button.
func WelcomePage(w http.ResponseWriter, r *http.Request) {
	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	state = randToken()
	session, _ := store.Get(r, "sess")
	session.Values["state"] = state
	session.Save(r, w)
	w.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + oauthCfg.AuthCodeURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}

// AuthCallback is where the user is directed to after logging in with Google.
func AuthCallback(w http.ResponseWriter, r *http.Request) {
	// restore the session and check for a match.
	session, err := store.Get(r, "sess")
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		fmt.Fprintln(w, "no state match; possible csrf OR cookies not enabled")
		return
	}

	// https://godoc.org/golang.org/x/oauth2#example-Config
	// Use the authorization code that is pushed to the redirect URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by conf.Client will refresh the token as necessary.
	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		fmt.Fprintln(w, "there was an issue getting your token")
		return
	}

	if !tkn.Valid() {
		fmt.Fprintln(w, "retreived invalid token")
		return
	}

	/*
		2018/01/27 17:58:49 Token :  {"access_token":"ya29.GlxPBU3mtkuc738UlO_40zAfperHAnQEn_ushyoBFTsiETy1uIuznIuE08nTgNKw2oEQ4iuPhRgG4u1W0uyoUvaYHUI8RPy3lEG-TVE3TpqlvQpso-k5Z8t8dS15sA","token_type":"Bearer","expiry":"2018-01-27T18:58:49.7588939-05:00"}

		token {
			access_token	A token that can be sent to a Google API.
			id_token	A JWT that contains identity information about the user that is digitally signed by Google.
			expires_in	The remaining lifetime of the access token.
			token_type	Identifies the type of token returned. At this time, this field always has the value Bearer.
			refresh_token (optional)	This field is only present if access_type=offline is included in the authentication request. For details, see Refresh tokens.
		}
	*/

	/* *********************************************************************************************************
	*
	*   START TEST CODE to see if the verify token code works.
	*
	 */

	pprint(tkn, "Token")

	// get the raw token.
	rawIDToken, ok := tkn.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	verifyToken(rawIDToken)

	/*
	*
	*   END TEST CODE to see if the verify token code works.
	*
	* ******************************************************************************************************** */

	// Now we have the user's token, we can create a client to hit the Google API we want.
	client := oauthCfg.Client(oauth2.NoContext, tkn)

	// get the data for the scope we requested.
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println("error getting userinfo")
		return
	}

	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	log.Println("User Info Profile: ", string(data))
	var result UserInfo
	json.Unmarshal(data, &result)

	// just for fun, insert them in People.

	person := Person{}
	person.ID = strconv.Itoa(len(people) + 1)
	person.Firstname = result.GivenName
	person.Lastname = result.FamilyName
	people = append(people, person)

	session.Values["email"] = result.Email
	session.Values["accessToken"] = tkn.AccessToken
	session.Save(r, w)

	// redirect to the page after they have authed.
	http.Redirect(w, r, "/people", 302)

}

func tokenSignIn(w http.ResponseWriter, r *http.Request) {

	// get the token.
	// rawIDToken, ok := tkn.Extra("id_token").(string)
	// if !ok {
	// 	http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
	// 	return
	// }

	// verifyToken(rawIDToken)

	// After obtaining user information from the ID token, you should query your app's user database. If the user already exists in your database, you should start an application session for that user.
	// If the user does not exist in your user database, you should redirect the user to your new-user sign-up flow. You may be able to auto-register the user based on the information you receive from Google,
	// or at the very least you may be able to pre-populate many of the fields that you require on your registration form. In addition to the information in the ID token, you can get additional user profile
	// information at our user profile endpoints.
	// https://developers.google.com/identity/protocols/OpenIDConnect

	// so we should create a new session if needed.
}

// GetPeople displays all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// GetPerson displays a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// DeletePerson deletes an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
