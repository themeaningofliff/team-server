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

func dbGetPlayerByID(id int) (Player, error) {
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
	err = db.QueryRow("SELECT * FROM players WHERE id=$1", id).Scan(&result.ID, &result.Firstname, &result.Lastname, &result.Email, &result.Phone, &result.Gender, &result.Address.Zipcode, &result.Active, &result.SignedUp, &result.CreatedOn)
	if err == sql.ErrNoRows {
		// no data, execute insert script to insert default data.
		log.Printf("No Players? %v\n", err)
		return result, ErrNoPlayer
	} else if err != nil {
		log.Printf("Error: %v\n", err)
		return result, ErrUnknown
	}

	return result, nil
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
	err = db.QueryRow("SELECT * FROM players WHERE email=$1", email).Scan(&result.ID, &result.Firstname, &result.Lastname, &result.Email, &result.Phone, &result.Gender, &result.Address.Zipcode, &result.Active, &result.SignedUp, &result.CreatedOn)
	if err == sql.ErrNoRows {
		// no data, execute insert script to insert default data.
		log.Printf("No Players? %v\n", err)
		return result, ErrNoPlayer
	} else if err != nil {
		log.Printf("Error: %v\n", err)
		return result, ErrUnknown
	}

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
	defer rows.Close()
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

		result = append(result, *np)
	}

	return result, nil
}

func dbInsertPlayer(ip Player) (int, error) {
	// connect to the DB
	var err error
	db, err := sql.Open("postgres", DBConnectionString)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		return -1, ErrDbNoConnection
	}
	defer db.Close()

	// check if there is any data
	var newID int
	err = db.QueryRow("INSERT INTO players(first_name,last_name,email,phone,gender,zipcode,active,signed_up) VALUES($1,$2,$3,$4,$5,$6,$7,$8) returning id",
		ip.Firstname, ip.Lastname, ip.Email, ip.Phone, ip.Gender, ip.Address.Zipcode, ip.Active, ip.SignedUp).Scan(&newID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return -1, ErrUnknown
	}

	ip.ID = newID
	log.Printf("Inserted new player: %s\n", ip.ToString())

	return newID, nil
}

func dbUpdatePlayer(ip Player) error {
	// connect to the DB
	var err error
	db, err := sql.Open("postgres", DBConnectionString)
	if err != nil {
		log.Printf("DB Init error: %v\n", err)
		return ErrDbNoConnection
	}
	defer db.Close()

	sqlStatement := `
	UPDATE players
	SET first_name=$2, last_name=$3, phone=$4, gender=$5, zipcode=$6, active=$7, signed_up=$8
	WHERE id = $1;`
	res, err := db.Exec(sqlStatement, ip.ID, ip.Firstname, ip.Lastname, ip.Phone, ip.Gender, ip.Address.Zipcode, ip.Active, ip.SignedUp)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return ErrUnknown
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return ErrUnknown
	}
	fmt.Printf("Updated %d players\n", affect)

	return nil
}

// TODO FOR LATER: Deleting rows.
// It's possible to prepare DB statements if you want to run multiple at the same time.
// https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.4.html

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
// TODO FOR LATER: Deleting rows.
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
