package matchmaking

import (
	"errors"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
)

var ErrPlayerIdEmpty = errors.New("player ID cannot be empty")
var ErrPlayerNotFound = errors.New("player not found")

var waitingCompetitions = make([]*model.Competition, 10)
var waitingPlayers = make([]*model.CompetingPlayer, 10)

func JoinCompetition(playerID string) (*model.Competition, error) {
	if playerID == "" {
		return nil, ErrPlayerIdEmpty
	}

	player := storage.Players[playerID]
	if player == nil {
		return nil, ErrPlayerNotFound
	}

	return nil, nil
}
