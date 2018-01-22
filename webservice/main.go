package webservice

import (
	// "log"
	"net/http"
)

// main function to boot up everything
func init() {
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})

    var router = NewRouter()

	// The path "/" matches everything not matched by some other path
	// in this case, redirect everything to our router.
	http.Handle("/", router)
	
	// Don't listen when running with Google App Engine
	// log.Fatal(http.ListenAndServe(":8000", router))
}
