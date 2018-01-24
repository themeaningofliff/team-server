CREATE DATABASE team_server;

CREATE TABLE IF NOT EXISTS players (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20), -- check max phone
    gender VARCHAR(20), -- simple or complex? 
    zipcode VARCHAR(9),-- consider international
    active BOOLEAN,
    signed_up BOOLEAN,
		UNIQUE (email, phone)
);

CREATE TABLE IF NOT EXISTS gameDefinition (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(20) NOT NULL,
		type VARCHAR(20) NOT NULL, -- eventually id to table if we want
);

CREATE TABLE IF NOT EXISTS playerGame (
		player_id INT NOT NULL,
		gameDefinition_id INT NOT NULL,
		base_skill_level VARCHAR(20), -- beginner, int, adv, pro
		custom_skill_level VARCHAR(20), -- sport specific string
		years_played VARCHAR(20), -- ranges
		frequency VARCHAR(20),
		allow_player_matching BOOLEAN
);

CREATE TABLE IF NOT EXISTS event ( -- event is a 
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		gameDefinition_id INT NOT NULL,
		created_on DATETIME DEFAULT NOW(),
		event_date DATETIME
		-- event time, start and end. event type (practice, match, etc)
);

CREATE TABLE IF NOT EXISTS gamePlayer (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		game_id INT NOT NULL,
		player_id INT NOT NULL,
		score INT
		-- start and end times
);