package webservice

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	oauth2 "golang.org/x/oauth2"
	// Google API Client

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var state string
var store = sessions.NewCookieStore([]byte("secret")) // TODO: This probably needs to be changed.

var ExUnknown = Exception{Code: 1, Message: "Unknown Error"}
var ExInvalid = Exception{Code: 2, Message: "Invalid request"}
var ExMandatory = Exception{Code: 3, Message: "Missing mandatory data"}
var ExForbidden = Exception{Code: 4, Message: "Forbidden"}
var ExBadRequest = Exception{Code: 5, Message: "Bad Request"}

// httpError generic error response wrapper
func httpError(w http.ResponseWriter, code int, ex Exception) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ex)
}

// httpUnknownError indicates server side error.
func httpBadRequest(w http.ResponseWriter, err error) {
	log.Printf("Bad request %v\n", err)
	httpError(w, http.StatusBadRequest, ExBadRequest)
}

// httpUnknownError indicates server side error.
func httpUnknownError(w http.ResponseWriter, err error) {
	log.Printf("Unknown error %v\n", err)
	httpError(w, http.StatusInternalServerError, ExUnknown)
}

// httpForbidden should be used to indicate access denied.
func httpForbidden(w http.ResponseWriter) {
	httpError(w, http.StatusForbidden, ExForbidden)
}

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
		httpUnknownError(w, err)
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

		// TODO: We need to figure out if we will force them to use the same email from the token?
		player, err := dbGetPlayerByEmail(tokenInfo.Email)
		existingUser := err == nil // if there has been no error, then we have found an entry.

		if existingUser {
			// the user has an account already, redirect to players.
			pprint(tokenInfo, "Found existing User")

			// TODO: What if the user hasn't signed up or is inactive?

			// TODO: This won't work - https://stackoverflow.com/questions/36345696/golang-http-redirect-with-headers
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", rawIDToken))
			http.Redirect(w, r, "/players", 302)
			return
		}

		if err != ErrNoPlayer {
			// unknown error!
			httpUnknownError(w, err)
			return
		}

		// the user does not have an account,
		// we should redirect them to a create profile page that is pre-populated with details we have got from Google.

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

		if existingUser && !player.SignedUp {
			log.Println("Non-Signed Up Player is Signing Up! Hooray, crack open the beers!")
			// TODO: Populate extra information from the registered player?
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

		// session.Values["email"] = userProfile.Email
		// session.Values["accessToken"] = tkn.AccessToken
		// session.Save(r, w)

	} else {
		pprint(tokenInfo, "Unverified email in Token")
		http.Error(w, "Unverified Email", http.StatusInternalServerError)
		return
	}

}

// GetPlayers displays all from the players var
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	players, err := dbGetPlayers()
	if err != nil {
		httpUnknownError(w, err)
		return
	}

	json.NewEncoder(w).Encode(players)
}

// GetPlayer displays a single data
func GetPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	players, err := dbGetPlayers()
	if err != nil {
		httpUnknownError(w, err)
		return
	}

	params := mux.Vars(r)
	for _, item := range players {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Player{})
}

func VerifyPlayerEmail(w http.ResponseWriter, r *http.Request, playerEmail string) bool {
	dummyToken, ok := context.Get(r, "dummy_token").(bool)
	if !ok {
		dummyToken = false
	}

	tokenEmail, ok := context.Get(r, "token_email").(string)
	if ok {
		if !(dummyToken || strings.EqualFold(tokenEmail, playerEmail)) {
			// email mismatch. Possibly someone is trying to hack an account
			log.Printf("Email Conflict CreatePlayer: %s vs %s\n", tokenEmail, playerEmail)
			httpError(w, http.StatusForbidden, Exception{Code: 901, Message: "User has already signed up and is active"})
			return false
		}
	} else {
		/* not string or not specified */
		// respond with access denied here.
		httpForbidden(w)
		return false
	}

	return true
}

// CreatePlayer creates a new item
func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		httpBadRequest(w, err)
		return
	}

	// validate the request.
	// TODO: Implement using some validation framework.
	if player.Email == "" || player.Phone == "" || player.Firstname == "" || player.Lastname == "" {
		httpError(w, http.StatusUnprocessableEntity, ExMandatory)
		return
	}

	if !(player.Active && player.SignedUp) {
		httpError(w, http.StatusUnprocessableEntity, Exception{Code: 100, Message: "New players should be signed up and active"})
		return
	}

	if !VerifyPlayerEmail(w, r, player.Email) {
		return
	}

	// check to see if the player is already there.
	existingPlayer, err := dbGetPlayerByEmail(player.Email)
	if err != nil && err != ErrNoPlayer {
		// some other error went down.
		httpUnknownError(w, err)
		return
	}

	if err == ErrNoPlayer {
		// expected error indicating there is no player.
		log.Println("Creating Player with email " + player.Email)

		newid, err := dbInsertPlayer(player)
		if err != nil {
			// some other error went down.
			httpUnknownError(w, err)
			return
		}

		player.ID = newid
		log.Printf("Inserted new player: %s\n", player.ToString())

		// return the newly created player, or should we just return their ID?
		json.NewEncoder(w).Encode(player)

	} else {
		// at this point there is no error indicating we did find a person row
		log.Printf("Person already exists: %s\n", existingPlayer.ToString())

		if existingPlayer.SignedUp {
			if existingPlayer.Active {
				// if the user has signed up and is active, return a failure indicating the user is already signed up.
				httpError(w, http.StatusConflict, Exception{Code: 101, Message: "User has already signed up and is active"})
				return
			}

			// user was signed up but is no longer active. Should we ask them to reactivate?
			httpError(w, http.StatusConflict, Exception{Code: 102, Message: "User was signed up but is no longer active. Reactive?"})
			return
		}

		// if the user hasn't signed up yet, update the existing row with these details
		// TODO: Audit?
		player.ID = existingPlayer.ID
		player.CreatedOn = existingPlayer.CreatedOn

		err := dbUpdatePlayer(player)
		if err != nil {
			// some other error went down.
			httpUnknownError(w, err)
			return
		}

		// successfully updated! Return the newly created player, or should we just return their ID?
		log.Printf("Updated new player: %s\n", player.ToString())
		json.NewEncoder(w).Encode(player)
	}
}

// DeletePlayer deactivates a player.
func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// parse the id from the url
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		httpBadRequest(w, err)
		return
	}

	// check to see if the player is already there.
	existingPlayer, err := dbGetPlayerByID(id)
	if err != nil {
		if err == ErrNoPlayer {
			// expected error indicating there is no player.
			// TODO: Security through obscurity - we shouldn't return a different error message.
			httpError(w, http.StatusUnprocessableEntity, Exception{Code: 105, Message: "Invalid player"})
			return
		}

		// some other error went down.
		httpUnknownError(w, err)
		return
	}

	if !VerifyPlayerEmail(w, r, existingPlayer.Email) {
		return
	}

	// deactivate player.
	existingPlayer.Active = false

	// TODO: Audit
	err = dbUpdatePlayer(existingPlayer)
	if err != nil {
		// some other error went down.
		httpUnknownError(w, err)
		return
	}

	// successfully updated! Return the newly created player, or should we just return their ID?
	log.Printf("Deactivated player: %s\n", existingPlayer.ToString())
	json.NewEncoder(w).Encode(existingPlayer)
}
