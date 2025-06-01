package leaderboard

import (
	"errors"
	"leaderboard/internal/storage"
	"leaderboard/internal/timeprovider"
)

var (
	ErrPlayerIdEmpty          = errors.New("player ID cannot be empty")
	ErrPlayerNotFound         = errors.New("player not found")
	ErrCompetitionEnded       = errors.New("competition has ended, cannot add score for player")
	ErrCompetitionNotStarted  = errors.New("competition has not started yet, cannot add score for player")
	ErrPlayerNotInCompetition = errors.New("player is not in a competition, cannot add score")
)

var AddScore = func(playerID string, points int) error {
	if playerID == "" {
		return ErrPlayerIdEmpty
	}
	player, playerFound := storage.Players[playerID]

	if !playerFound {
		return ErrPlayerNotFound
	}
	comp := player.Competition()
	if comp != nil {
		// TODO: This is domain logic. Move to the model
		if comp.StartedAt().IsZero() {
			return ErrCompetitionNotStarted
		} else if comp.EndsAt().Before(timeprovider.Current.Now()) {
			return ErrCompetitionEnded
		}
	} else {
		return ErrPlayerNotInCompetition
	}

	err := comp.AddScore(player, points)
	return err
}
