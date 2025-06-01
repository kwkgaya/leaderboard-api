package config

import "time"

// TODO: Read from config or env

var (
	MatchWaitDuration  = 30 * time.Second
	MatchRetryInterval = 1 * time.Second
)

const (
	MaxPlayersForCompetition int = 10
	MinPlayersForCompetition int = 2
	CompetitionDuration          = 1 * time.Hour

	MaxLevel = 10 // Maximum level a player can have
	MinLevel = 1  // Minimum level a player can have
)
