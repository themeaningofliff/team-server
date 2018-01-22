package main

// Player Type (more like an object)
type Player struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"first_name,omitempty"`
	Lastname  string   `json:"last_name,omitempty"`
	Email     string   `json:"email,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Address   *Address `json:"address,omitempty"`
	// CreatedOn Time.time `json:"created_on,omitempty"`
	CreatedOn string `json:"created_on,omitempty"`
	Active    bool   `json:"active,omitempty"`
	SignedUp  bool   `json:"signed_up,omitempty"`
}

// Address Type
type Address struct {
	ID      int    `json:"id,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zipcode string `json:"zip_code,omitempty"`
}

// Players - group of players
var players []Player

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

// PlayerGames - group of games a player plays
var playerGames []PlayerGame

// Game - game played with details
type Game struct {
	ID               int `json:"id,omitempty"`
	GameDefinitionID int `json:"game_definition_id,omitempty"`
	// GameStarted      Time.time `json:"game_started,omitempty"`
	// CreatedOn        Time.time `json:"created_on,omitempty"`
	GameStarted string `json:"game_started,omitempty"`
	CreatedOn   string `json:"created_on,omitempty"`
}

// Games - group of games
var games []Game

// GameScore - scores of players by game
type GameScore struct {
	ID       int `json:"id,omitempty"`
	GameID   int `json:"game_id,omitempty"`
	PlayerID int `json:"player_id,omitempty"`
	Score    int `json:"score,omitempty"`
}
