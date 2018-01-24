package webservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	// "google.golang.org/appengine"
	// "google.golang.org/appengine/user"

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

func getLoginURL(state string) string {
	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return oauthCfg.AuthCodeURL(state)
}

// Welcome greets a user who has signed in to the app with a personalized message and a link to sign out. If the user is not signed in, the app offers a link to the sign-in page for Google Accounts.
func Welcome(w http.ResponseWriter, r *http.Request) {
	state = randToken()
	session, _ := store.Get(r, "sess")
	session.Values["state"] = state
	session.Save(r, w)
	w.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))

	// w.Header().Set("Content-type", "text/html; charset=utf-8")
	// ctx := appengine.NewContext(r)
	// u := user.Current(ctx)
	// if u == nil {
	// 	url, _ := user.LoginURL(ctx, "/")
	// 	fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
	// 	return
	// }
	// url, _ := user.LogoutURL(ctx, "/")
	// fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

// func WelcomeOAuth(w http.ResponseWriter, r *http.Request) {
// 	ctx := appengine.NewContext(r)
// 	u, err := user.CurrentOAuth(ctx, "")
// 	if err != nil {
// 		http.Error(w, "OAuth Authorization header required", http.StatusUnauthorized)
// 		return
// 	}
// 	if !u.Admin {
// 		http.Error(w, "Admin login only", http.StatusUnauthorized)
// 		return
// 	}
// 	fmt.Fprintf(w, `Welcome, admin user %s!`, u)
// }

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sess")
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	// // Retrieve our struct and type-assert it
	// val := session.Values["state"]
	// var st string
	// if _, ok := val.(string); !ok {
	// 	log.Printf("Could not read state from session")
	// }
	// log.Println("state : ", st)

	if r.URL.Query().Get("state") != session.Values["state"] {
		fmt.Fprintln(w, "no state match; possible csrf OR cookies not enabled")
		return
	}

	// https://godoc.org/golang.org/x/oauth2#example-Config
	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		fmt.Fprintln(w, "there was an issue getting your token")
		return
	}

	if !tkn.Valid() {
		fmt.Fprintln(w, "retreived invalid token")
		return
	}

	client := oauthCfg.Client(oauth2.NoContext, tkn)

	// get the data for the scope we requested.
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		fmt.Println(w, "error getting userinfo")
		return
	}

	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	log.Println("Resp body: ", string(data))
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
