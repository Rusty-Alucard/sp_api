SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;

/**************
* DDL
**************/

DROP TABLE IF EXISTS `events`;

CREATE TABLE `events` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`name` VARCHAR(255) NOT NULL,
`name_jp` VARCHAR(255) NOT NULL,
PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `teams`;

CREATE TABLE `teams` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`fifa_trigramma` char(3) NOT NULL,
`name` VARCHAR(255) NOT NULL,
`name_jp` VARCHAR(255) NOT NULL,
`confederation` VARCHAR(255) NOT NULL,
PRIMARY KEY(`id`),
KEY `index_teams_on_fifa_trigramma` (`fifa_trigramma`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `fifa_ranking`;

CREATE TABLE `fifa_ranking` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`fifa_trigramma` char(3) NOT NULL,
`date` date,
`rank` int(3) unsigned NOT NULL,
`points` int(5) unsigned NOT NULL,
PRIMARY KEY(`id`),
KEY `index_teams_on_fifa_trigramma` (`fifa_trigramma`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `coaches`;

CREATE TABLE `coaches` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`name` VARCHAR(255) NOT NULL,
`name_jp` VARCHAR(255) NOT NULL,
`birthday` date,
PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `participants`;

CREATE TABLE `participants` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`event_id` bigint(20) unsigned NOT NULL,
`team_id` bigint(20) unsigned NOT NULL,
`coach_id` bigint(20) unsigned,
PRIMARY KEY(`id`),
KEY `index_participants_on_event_id` (`event_id`),
KEY `index_participants_on_team_id` (`team_id`),
KEY `index_participants_on_coach_id` (`coach_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `players`;

CREATE TABLE `players` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`name` VARCHAR(255) NOT NULL,
`name_jp` VARCHAR(255) NOT NULL,
`nickname` VARCHAR(255),
`birthday` date,
PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `registered_players`;

CREATE TABLE `registered_players` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`participant_id` bigint(20) unsigned NOT NULL,
`player_id` bigint(20) unsigned NOT NULL,
`number` int(2) unsigned,
`captain` boolean NOT NULL DEFAULT false,
`key_player` boolean NOT NULL DEFAULT false,
PRIMARY KEY(`id`),
KEY `index_registered_players_on_participant_id` (`participant_id`),
KEY `index_registered_players_on_player_id` (`player_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `stadiums`;

CREATE TABLE `stadiums` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`name` VARCHAR(255) NOT NULL,
`name_jp` VARCHAR(255) NOT NULL,
`venue` VARCHAR(255) NOT NULL,
`venue_jp` VARCHAR(255) NOT NULL,
`timezone` VARCHAR(255) NOT NULL,
PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;


DROP TABLE IF EXISTS `matches`;

CREATE TABLE `matches` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`event_id` bigint(20) unsigned NOT NULL,
`studium_id` bigint(20) unsigned NOT NULL,
`kickoff` timestamp,
`groups` varchar(255),
PRIMARY KEY(`id`),
KEY `index_matches_on_event_id` (`event_id`),
KEY `index_matches_on_studium_id` (`studium_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `starting_lineups`;

CREATE TABLE `starting_lineups` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`match_id` bigint(20) unsigned NOT NULL,
`registered_player_id` bigint(20) unsigned NOT NULL,
`position_x` int(4),
`position_y` int(4),
PRIMARY KEY(`id`),
KEY `index_starting_lineups_on_match_id` (`match_id`),
KEY `index_starting_lineups_on_registered_player_id` (`registered_player_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `match_progresses`;

CREATE TABLE `match_progresses` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`match_id` bigint(20) unsigned NOT NULL,
`registered_player_id` bigint(20) unsigned NOT NULL,
`type` int(2) unsigned NOT NULL,
`time` int(3) unsigned NOT NULL,
`additional_time` int(3) unsigned NOT NULL DEFAULT 0,
`action` int(2) unsigned NOT NULL,
PRIMARY KEY(`id`),
KEY `index_match_progresses_on_match_id` (`match_id`),
KEY `index_match_progresses_on_registered_player_id` (`registered_player_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `opponents`;

CREATE TABLE `opponents` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`match_id` bigint(20) unsigned NOT NULL,
`participant_id` bigint(20) unsigned NOT NULL,
`home` boolean NOT NULL DEFAULT false,
PRIMARY KEY(`id`),
KEY `index_opponents_on_match_id` (`match_id`),
KEY `index_opponents_on_registered_player_id` (`participant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;

