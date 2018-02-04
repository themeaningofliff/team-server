package webservice

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	oauth2 "golang.org/x/oauth2"
	// Google API Client

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

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
	w.Write([]byte("<html><title>Welcome</title> <body><H2>Welcome!</H2><BR><a href='" + oauthCfg.AuthCodeURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}

// WelcomePage greets a user who has found the site and presents them with a Login with Google button.
func LoginAgain(w http.ResponseWriter, r *http.Request) {
	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	state = randToken()
	session, _ := store.Get(r, "sess")
	session.Values["state"] = state
	session.Save(r, w)
	w.Write([]byte("<html><title>Login</title> <body><H2>You've been logged off. Please login again</H2><BR><a href='" + oauthCfg.AuthCodeURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
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

	// TODO: Read this later - https://tools.ietf.org/html/rfc6819
	if !tkn.Valid() {
		fmt.Fprintln(w, "retreived invalid token")
		return
	}

	pprint(tkn, "Token")

	/*
		token {
			access_token	A token that can be sent to a Google API.
			id_token	A JWT that contains identity information about the user that is digitally signed by Google.
			expires_in	The remaining lifetime of the access token.
			token_type	Identifies the type of token returned. At this time, this field always has the value Bearer.
			refresh_token (optional)	This field is only present if access_type=offline is included in the authentication request. For details, see Refresh tokens.
		}

		Example:
		018/01/31 07:06:17 Token
		{
			"access_token": "ya29.GlxTBRS9CNhLTCRIp3qgTZ9uwSOYGxAVo09qjFB2b8rdtfcjPGQmmHWTr1zae83YM3vy4HrITafvdxt5dXT2eXeYFE8bP2vzz0cLGa9fyR9SfafQSeuOqAOv2Q-mCw",
			"token_type": "Bearer",
			"expiry": "2018-01-31T08:06:17.2414786-05:00"
		}

	*/

	// After obtaining user information from the ID token, you should query your app's user database. If the user already exists in your database, you should start an application session for that user.
	// If the user does not exist in your user database, you should redirect the user to your new-user sign-up flow. You may be able to auto-register the user based on the information you receive from Google,
	// or at the very least you may be able to pre-populate many of the fields that you require on your registration form. In addition to the information in the ID token, you can get additional user profile
	// information at our user profile endpoints.
	// https://developers.google.com/identity/protocols/OpenIDConnect

	// get the raw OpenID token.
	rawIDToken, ok := tkn.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	pprint(rawIDToken, "ID Token")

	// verify & parse the token
	tokenInfo, err := verifyTokenAPI(rawIDToken)
	if err != nil {
		log.Println("Invalid authorization token: " + err.Error())
		http.Error(w, "Invalid authorization token", http.StatusInternalServerError)
		return
	}

	if tokenInfo.VerifiedEmail {
		var existingUser = PersonAlreadyExistsByEmail(tokenInfo.Email)

		if !existingUser {
			// the user does not have an account, we should redirect them to a create profile page
			// that is pre-populated with details we have got from Google.

			// Now we have the user's token, we can create a client to hit the Google API we want.
			client := oauthCfg.Client(oauth2.NoContext, tkn)

			// get the data for the scope we requested - in this case, Google Profile UserInfo
			userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
			if err != nil {
				log.Println("error getting userinfo")
				return
			}

			defer userinfo.Body.Close()
			data, _ := ioutil.ReadAll(userinfo.Body)
			log.Println("User Info Profile: ", string(data))
			var userProfile UserInfo
			json.Unmarshal(data, &userProfile)

			// The problem with this token is we have to pass it in plain text in the html
			// TODO: Look into http://www.gorillatoolkit.org/pkg/csrf or https://github.com/justinas/nosurf
			var f interface{}
			f = map[string]interface{}{
				"Firstname":  userProfile.GivenName,
				"Lastname":   userProfile.FamilyName,
				"Email":      userProfile.Email,
				"tokenField": rawIDToken,
			}
			t, err := template.ParseFiles("./webservice/createProfile.html")
			if err != nil {
				log.Println("error parsing template " + err.Error())
				return
			}
			err = t.Execute(w, f)
			if err != nil {
				log.Println("error executing template " + err.Error())
				return
			}

			// To test that the above "is user" check works, insert them in People so next time we refresh it doesn't do this.
			// person := Person{}
			// person.ID = strconv.Itoa(len(people) + 1)
			// person.Email = userProfile.Email
			// person.Firstname = userProfile.GivenName
			// person.Lastname = userProfile.FamilyName
			// people = append(people, person)

			// session.Values["email"] = userProfile.Email
			// session.Values["accessToken"] = tkn.AccessToken
			// session.Save(r, w)

		} else {
			// the user has an account already, redirect to people.
			pprint(tokenInfo, "Found existing User")
			w.Header().Set("Authorization", "Bearer "+rawIDToken)
			http.Redirect(w, r, "/people", 302)
		}
	} else {
		pprint(tokenInfo, "Unverified email in Token")
		http.Error(w, "Unverified Email", http.StatusInternalServerError)
		return
	}

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
	// params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	if !PersonAlreadyExistsByEmail(person.Email) {
		log.Println("Creating Person with email " + person.Email)
		person.ID = strconv.Itoa(len(people) + 1)
		people = append(people, person)
	} else {
		log.Println("Person already exists with email " + person.Email)
	}

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
