CREATE DATABASE IF NOT EXISTS team_server;

CREATE TABLE IF NOT EXISTS players (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20), #check max phone
    gender VARCHAR(20), #simple or complex? 
    zipcode VARCHAR(9) #consider international
)

CREATE TABLE IF NOT EXISTS gameDefinition (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(20),
		type VARCHAR(20), #eventually id to table if we want
)

CREATE TABLE IF NOT EXISTS playerGame (
		player_id INT,
		gameDefinition_id INT,
		base_skill_level VARCHAR(20), #beginner, int, adv, pro
		custom_skill_level VARCHAR(20), # sport specific string
		years_played VARCHAR(20), #ranges
		allow_player_matching BOOLEAN 
)

CREATE TABLE IF NOT EXISTS game (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		gameDefinition_id INT,
		created_on DATETIME,
		game_started DATETIME
)

CREATE TABLE IF NOT EXISTS gameScore (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		game_id INT,
		player_id INT,
		score INT
)