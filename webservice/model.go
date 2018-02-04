package webservice

import "strings"

// Exception wraps a json error message
type Exception struct {
	Message string `json:"message"`
}

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

// Person Type (more like an object)
type Person struct {
	ID        string   `json:"id,omitempty"`
	Email     string   `json:"email,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// Address Type
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

// PersonAlreadyExistsByEmail checks if a person already exists by email
func PersonAlreadyExistsByEmail(email string) bool {
	// check if the user has an account already.
	for _, item := range people {
		// We could actually use Google's UserIDs if we wanted but might not be good if we use another identify provider.
		// tokenInfo.UserId // "user_id": "100682826382643775970",

		if strings.EqualFold(item.Email, email) { //case insensitve comparison
			return true
		}
	}

	return false
}

var people []Person
