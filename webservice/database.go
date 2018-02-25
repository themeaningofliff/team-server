package webservice

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const DBConnectionString = "postgres://pgteam:pgteam1234@localhost/team_server?sslmode=disable"

var ErrDbNoConnection = errors.New("Failed to connect to DB")
var ErrDBFailure = errors.New("Database error")
var ErrNoPlayers = errors.New("No Players found")
var ErrNoPlayer = errors.New("Player not found")
var ErrUnknown = errors.New("Unknown error")

// executeDBScript runs the sql passed.
func executeDBScript(db *sql.DB, inputFile string) {
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
	// connect to the DB
	db, err := sql.Open("postgres", DBConnectionString)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		panic(err)
	}
	defer db.Close()

	// check if the connection is valid
	if err = db.Ping(); err != nil {
		// if there is no database, create it.
		// TODO: Remove this and put into some utility tool/script.
		if strings.Contains(err.Error(), "does not exist") {
			log.Println("team_server db does not exist!!")
			executeDBScript(db, "database/database.sql")
		} else {
			// some other error, perhaps permissions?
			log.Fatal(err)
			return
		}
	}

	// check if there is any data
	_, err = db.Query("SELECT * from players")
	if err == sql.ErrNoRows {
		// no data, execute insert script to insert default data.
		log.Printf("No Rows: %v\n", err)
		executeDBScript(db, "database/test_data.sql")
		return
	} else if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
}

func dbGetPlayerByEmail(email string) (Player, error) {
	var result Player

	// connect to the DB
	var err error
	db, err := sql.Open("postgres", DBConnectionString)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		return result, ErrDbNoConnection
	}
	defer db.Close()

	// check if there is any data
	var np = new(Player)
	err = db.QueryRow("SELECT * FROM players WHERE email=$1", email).Scan(&np.ID, &np.Firstname, &np.Lastname, &np.Email, &np.Phone, &np.Gender, &np.Address.Zipcode, &np.Active, &np.SignedUp, &np.CreatedOn)
	if err == sql.ErrNoRows {
		// no data, execute insert script to insert default data.
		log.Printf("No Players? %v\n", err)
		return result, ErrNoPlayer
	} else if err != nil {
		log.Printf("Error: %v\n", err)
		return result, ErrUnknown
	}

	fmt.Println(np.ToString())

	return result, nil
}

func dbGetPlayers() ([]Player, error) {
	var result []Player // empty.

	// connect to the DB
	var err error
	db, err := sql.Open("postgres", DBConnectionString)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		return result, ErrDbNoConnection
	}
	defer db.Close()

	// check if there is any data
	rows, err := db.Query("SELECT * from players")
	if err == sql.ErrNoRows {
		// no data, execute insert script to insert default data.
		log.Printf("No Players? %v\n", err)
		return result, ErrNoPlayers
	} else if err != nil {
		log.Printf("Error: %v\n", err)
		return result, ErrUnknown
	}

	// read in the current list of players.
	// TODO: This is not required in an initialization.
	for rows.Next() {
		var np = new(Player)
		// np.Address = new(Address)
		err = rows.Scan(&np.ID, &np.Firstname, &np.Lastname, &np.Email, &np.Phone, &np.Gender, &np.Address.Zipcode, &np.Active, &np.SignedUp, &np.CreatedOn)
		if err != nil {
			log.Printf("DB Init error: %v\n", err)
			return result, ErrDBFailure
		}
		fmt.Println(np.ToString())
		result = append(result, *np)
	}

	return result, nil

	// fmt.Println("# Inserting values")
	// var lastInsertId int
	// err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	// if err != nil {
	// 	log.Printf("DB Init error: %v\n", err)
	// 	panic(err)
	// }
	// fmt.Println("last inserted id =", lastInsertId)

	//
	//
	//

	// fmt.Println("# Updating")
	// stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
	// if err != nil {
	// 	log.Printf("DB Init error: %v\n", err)
	// 	panic(err)
	// }

	// res, err := stmt.Exec("astaxieupdate", lastInsertId)
	// if err != nil {
	// 	log.Printf("DB Init error: %v\n", err)
	// 	panic(err)
	// }

	// affect, err := res.RowsAffected()
	// if err != nil {
	// 	log.Printf("DB Init error: %v\n", err)
	// 	panic(err)
	// }
	// fmt.Println(affect, "rows changed")

	//
	//
	//

	// fmt.Println("# Deleting")
	// stmt, err = db.Prepare("delete from userinfo where uid=$1")
	// if err != nil {
	// 	log.Printf("DB Init error: %v\n", err)
	// 	panic(err)
	// }

	// res, err = stmt.Exec(lastInsertId)
	// if err != nil {
	// 	log.Printf("DB Init error: %v\n", err)
	// 	panic(err)
	// }

	// affect, err = res.RowsAffected()
	// if err != nil {
	// 	log.Printf("DB Init error: %v\n", err)
	// 	panic(err)
	// }

	// fmt.Println(affect, "rows changed")
}
