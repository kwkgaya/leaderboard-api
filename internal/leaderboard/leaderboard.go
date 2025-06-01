package leaderboard

import (
	"errors"
	"leaderboard/internal/model"
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

var AddScore = func(playerId string, points int) error {
	comp, err := getCompetition(playerId)
	if err != nil {
		return err
	}
	if comp.StartedAt().IsZero() {
		return ErrCompetitionNotStarted
	} else if comp.EndsAt().Before(timeprovider.Current.Now()) {
		return ErrCompetitionEnded
	}

	err = comp.AddScore(playerId, points)
	return err
}

func getCompetition(playerId string) (model.ICompetition, error) {
	if playerId == "" {
		return nil, ErrPlayerIdEmpty
	}
	player, found := storage.Players[playerId]
	if !found {
		return nil, ErrPlayerNotFound
	}
	comp := player.Competition()
	if comp == nil {
		return nil, ErrPlayerNotInCompetition
	}

	return comp, nil
}
