package webservice

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

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
var db *sql.DB

// executeDBScript runs the sql passed.
func executeDBScript(inputFile string) {
	file, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Printf("DB File %s error: %v\n", inputFile, err)
		return
	}

	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		log.Printf("DB request: %v\n", request)
		result, err := db.Exec(request)
		if err != nil {
			log.Printf("DB error: %v\n", err)
			return
		}
		log.Printf("DB result: %v\n", result)
	}
}

func initDB() {
	// init the DataBase
	conninfo := "postgres://pgteam:pgteam1234@localhost/team_server?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", conninfo)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			log.Println("team_server db does not exist!!")
			executeDBScript("database/database.sql")
		} else {
			log.Fatal(err)
			return
		}
	}

	// check if the database exists.
	var result int
	err = db.QueryRow("SELECT * from players").Scan(&result)
	if err == sql.ErrNoRows {
		log.Printf("No Rows: %v\n", err)
		executeDBScript("database/test_data.sql")
		return
	} else if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	log.Printf("Result: %v\n", result)

	fmt.Println("# Inserting values")

	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}
	fmt.Println("last inserted id =", lastInsertId)

	fmt.Println("# Updating")
	stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}

	res, err := stmt.Exec("astaxieupdate", lastInsertId)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}

	fmt.Println(affect, "rows changed")

	fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created time.Time
		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			log.Printf("DB Init error: %v\n", err)
			panic(err)
		}
		fmt.Println("uid | username | department | created ")
		fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
	}

	fmt.Println("# Deleting")
	stmt, err = db.Prepare("delete from userinfo where uid=$1")
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}

	res, err = stmt.Exec(lastInsertId)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}

	affect, err = res.RowsAffected()
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}

	fmt.Println(affect, "rows changed")
}

// main function to boot up everything
func init() {
	initDB()

	players = append(players, Player{ID: "1", Firstname: "John", Lastname: "Doe", Email: "john@funding.com", Phone: "5555555555", Address: &Address{ID: 1, City: "City X", State: "State X", Zipcode: "11101"}, CreatedOn: "1443492224", Active: true, SignedUp: true})
	players = append(players, Player{ID: "2", Firstname: "Koko", Lastname: "Doe", Email: "koko@funding.com", Phone: "5555551234", Address: &Address{ID: 2, City: "City Z", State: "State Y", Zipcode: "11101"}, CreatedOn: "1438947306", Active: false, SignedUp: false})

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
