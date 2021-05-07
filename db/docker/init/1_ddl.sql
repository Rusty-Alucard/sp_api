set client_encoding = 'UTF8';

\c sp

/*******
* DDL
*******/

DROP TABLE IF EXISTS events;

CREATE TABLE events (
id bigserial NOT NULL,
name VARCHAR(255) NOT NULL,
name_jp VARCHAR(255) NOT NULL,
PRIMARY KEY(id)
);

DROP TABLE IF EXISTS teams;

CREATE TABLE teams (
id bigserial NOT NULL,
fifa_trigramma char(3) NOT NULL,
name VARCHAR(255) NOT NULL,
name_jp VARCHAR(255) NOT NULL,
confederation VARCHAR(255) NOT NULL,
PRIMARY KEY(id)
);
CREATE INDEX index_teams_on_fifa_trigramma on teams (fifa_trigramma);

DROP TABLE IF EXISTS fifa_ranking;

CREATE TABLE fifa_ranking (
id bigserial NOT NULL,
fifa_trigramma char(3) NOT NULL,
date date,
rank smallint NOT NULL,
points smallint NOT NULL,
PRIMARY KEY(id)
);
CREATE INDEX index_fifa_ranking_on_fifa_trigramma on fifa_ranking (fifa_trigramma);

DROP TABLE IF EXISTS coaches;

CREATE TABLE coaches (
id bigserial NOT NULL,
name VARCHAR(255) NOT NULL,
name_jp VARCHAR(255) NOT NULL,
birthday date,
PRIMARY KEY(id)
);

DROP TABLE IF EXISTS participants;

CREATE TABLE participants (
id bigserial NOT NULL,
event_id integer NOT NULL,
team_id integer NOT NULL,
coach_id integer,
PRIMARY KEY(id)
);
CREATE INDEX index_participants_on_event_id on participants (event_id);
CREATE INDEX index_participants_on_team_id on participants (team_id);
CREATE INDEX index_participants_on_coach_id on participants (coach_id);

DROP TABLE IF EXISTS players;

CREATE TABLE players (
id bigserial NOT NULL,
name VARCHAR(255) NOT NULL,
name_jp VARCHAR(255) NOT NULL,
nickname VARCHAR(255),
birthday date,
PRIMARY KEY(id)
);

DROP TABLE IF EXISTS registered_players;

CREATE TABLE registered_players (
id bigserial NOT NULL,
participant_id integer NOT NULL,
player_id integer NOT NULL,
number smallint,
captain boolean NOT NULL DEFAULT false,
key_player boolean NOT NULL DEFAULT false,
PRIMARY KEY(id)
);
CREATE INDEX index_registered_players_on_participant_id on registered_players (participant_id);
CREATE INDEX index_registered_players_on_player_id on registered_players (player_id);

DROP TABLE IF EXISTS stadiums;

CREATE TABLE stadiums (
id bigserial NOT NULL,
name VARCHAR(255) NOT NULL,
name_jp VARCHAR(255) NOT NULL,
venue VARCHAR(255) NOT NULL,
venue_jp VARCHAR(255) NOT NULL,
timezone VARCHAR(255) NOT NULL,
PRIMARY KEY(id)
);


DROP TABLE IF EXISTS matches;

CREATE TABLE matches (
id bigserial NOT NULL,
event_id integer NOT NULL,
studium_id integer NOT NULL,
kickoff timestamp,
groups varchar(255),
PRIMARY KEY(id)
);
CREATE INDEX index_matches_on_event_id on matches (event_id);
CREATE INDEX index_matches_on_studium_id on matches (studium_id);

DROP TABLE IF EXISTS starting_lineups;

CREATE TABLE starting_lineups (
id bigserial NOT NULL,
match_id integer NOT NULL,
registered_player_id integer NOT NULL,
position_x smallint,
position_y smallint,
PRIMARY KEY(id)
);
CREATE INDEX index_starting_lineups_on_match_id on starting_lineups (match_id);
CREATE INDEX index_starting_lineups_on_registered_player_id on starting_lineups (registered_player_id);

DROP TABLE IF EXISTS match_progresses;

CREATE TABLE match_progresses (
id bigserial NOT NULL,
match_id integer NOT NULL,
registered_player_id integer NOT NULL,
type smallint NOT NULL,
time smallint NOT NULL,
additional_time smallint NOT NULL DEFAULT 0,
action smallint NOT NULL,
PRIMARY KEY(id)
);
CREATE INDEX index_match_progresses_on_match_id on match_progresses (match_id);
CREATE INDEX index_match_progresses_on_registered_player_id on match_progresses (registered_player_id);

DROP TABLE IF EXISTS opponents;

CREATE TABLE opponents (
id bigserial NOT NULL,
match_id integer NOT NULL,
participant_id integer NOT NULL,
home boolean NOT NULL DEFAULT false,
PRIMARY KEY(id)
);
CREATE INDEX index_opponents_on_match_id on opponents (match_id);
CREATE INDEX index_opponents_on_registered_player_id on opponents (participant_id);
