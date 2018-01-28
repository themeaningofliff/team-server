package webservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	oidc "github.com/coreos/go-oidc"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

var oauthCred Credentials
var oauthCfg *oauth2.Config
var oauthVerifier *oidc.IDTokenVerifier

// main function to boot up everything
func init() {
	// dummy data
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})

	// read the credentials.
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &oauthCred)

	/*
			  OAuth2 Client ID: 379625204959-4t2js39veijsiopjog6e2rtfruo0qrb3.apps.googleusercontent.com
		      OAuth2 Client Secret: rWJj9RaDvB7zUoYc3QSn8cPK
	*/
	// construct OAuth struct
	oauthCfg = &oauth2.Config{
		ClientID:     oauthCred.Cid,
		ClientSecret: oauthCred.Csecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/oauth2callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
	}

	// construct an oauth verifier for Google Accounts.
	// TODO: We are using OIDC, a non-Google API to do this. If Google ever releases one, we should use theirs.
	// https://developers.google.com/identity/sign-in/android/backend-auth
	provider, err := oidc.NewProvider(oauth2.NoContext, "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}
	oidcConfig := &oidc.Config{
		ClientID: oauthCred.Cid,
	}
	oauthVerifier = provider.Verifier(oidcConfig)

	// setup router.
	var router = NewRouter()

	// The path "/" matches everything not matched by some other path
	// in this case, redirect everything to our router.
	http.Handle("/", router)

	// Don't listen when running with Google App Engine
	log.Fatal(http.ListenAndServe(":8080", router))
}
