CREATE DATABASE team_server;

CREATE TABLE IF NOT EXISTS players (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20), -- check max phone
    gender VARCHAR(20), -- simple or complex? 
    zipcode VARCHAR(9),-- consider international
    active BOOLEAN NOT NULL DEFAULT false,
    signed_up BOOLEAN NOT NULL DEFAULT false,
	created_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),		
		UNIQUE (email, phone)
);

CREATE TABLE IF NOT EXISTS gameDefinition ( -- tennis, climbing etc.
		id SERIAL PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		type VARCHAR(20) NOT NULL -- eventually id to table if we want
);

CREATE TABLE IF NOT EXISTS playerGame ( -- which games the player is into
		player_id BIGINT REFERENCES players,
		gameDefinition_id INT REFERENCES gameDefinition,
		base_skill_level VARCHAR(20), -- beginner, int, adv, pro
		custom_skill_level VARCHAR(20), -- sport specific string
		years_played VARCHAR(20), -- ranges
		frequency VARCHAR(20),
		allow_player_matching BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS events ( -- event is a actual instance of a game/event
		id BIGSERIAL PRIMARY KEY,
		gameDefinition_id INT REFERENCES gameDefinition,
		created_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		event_date TIMESTAMP WITH TIME ZONE NOT NULL
		-- event time, start and end. event type (practice, match, etc)
);

CREATE TABLE IF NOT EXISTS eventPlayer (
		id BIGSERIAL PRIMARY KEY,
		event_id BIGINT REFERENCES events,
		player_id BIGINT REFERENCES players,
		score INT NOT NULL
		-- start and end times
);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pgteam;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO pgteam;
