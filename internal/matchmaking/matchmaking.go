package matchmaking

import (
	"fmt"
	"leaderboard/internal/storage"
	"leaderboard/internal/storage/model"
)

var ErrPlayerIdEmpty = fmt.Errorf("player ID cannot be empty")
var ErrPlayerNotFound = fmt.Errorf("player not found")

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
