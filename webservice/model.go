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

// Player Type (more like an object)
type Player struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"first_name,omitempty"`
	Lastname  string   `json:"last_name,omitempty"`
	Email     string   `json:"email,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Address   *Address `json:"address,omitempty"`
	CreatedOn string   `json:"created_on,string,omitempty"`
	Active    bool     `json:"active,omitempty"`
	SignedUp  bool     `json:"signed_up,omitempty"`
}

// Address Type
type Address struct {
	ID      int    `json:"id,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zipcode string `json:"zip_code,omitempty"`
}

// PlayerAlreadyExistsByEmail checks if a person already exists by email
func PlayerAlreadyExistsByEmail(email string) bool {
	// check if the user has an account already.
	for _, item := range players {
		// We could actually use Google's UserIDs if we wanted but might not be good if we use another identify provider.
		// tokenInfo.UserId // "user_id": "100682826382643775970",

		if strings.EqualFold(item.Email, email) { //case insensitve comparison
			return true
		}
	}

	return false
}

// Players - group of players
var players []Player

// PlayerGames - group of games a player plays
var playerGames []PlayerGame

// GameDefinition - Model for the base game
type GameDefinition struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// PlayerGame - mapping of players to games they play with details
type PlayerGame struct {
	PlayerID            int    `json:"id,omitempty"`
	GameDefinitionID    int    `json:"game_definition_id,omitempty"`
	BaseSkillLevel      string `json:"baseskilllevel,omitempty"`
	CustomSkillLevel    string `json:"custom_skill_level,omitempty"`
	YearsPlayed         string `json:"years_played,omitempty"`
	AllowPlayerMatching bool   `json:"allow_player_matching,omitempty"`
}

// Event - game played with details
type Event struct {
	ID               int    `json:"id,omitempty"`
	GameDefinitionID int    `json:"game_definition_id,omitempty"`
	GameStarted      string `json:"game_started,omitempty"`
	CreatedOn        string `json:"created_on,omitempty"`
}

// Games - group of games
var events []Event

// EventScore - scores of players by game
type EventScore struct {
	ID       int `json:"id,omitempty"`
	GameID   int `json:"game_id,omitempty"`
	PlayerID int `json:"player_id,omitempty"`
	Score    int `json:"score,omitempty"`
}
