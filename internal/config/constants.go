package config

import "time"

// TODO: Read from config or env

var (
	MatchWaitDuration       = 30 * time.Second
	MatchRetryInterval      = 1 * time.Second
	CompetitionDuration     = 1 * time.Hour
	MaxCompetitionsInMemory = 100
)

const (
	MaxPlayersForCompetition = 10
	MinPlayersForCompetition = 2

	MaxLevel = 10 // Maximum level a player can have
	MinLevel = 1  // Minimum level a player can have
)
