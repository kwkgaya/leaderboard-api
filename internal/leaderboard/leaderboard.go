package leaderboard

// TODO: Use interfaces instead of function varriables

import (
	"errors"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
	"leaderboard/internal/timeprovider"
	"time"
)

var (
	ErrPlayerIdEmpty          = errors.New("player ID cannot be empty")
	ErrLeaderboardIdEmpty     = errors.New("leaderboard ID cannot be empty")
	ErrPlayerNotFound         = errors.New("player not found")
	ErrCompetetionNotFound    = errors.New("competition not found")
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

var GetLeaderboard = func(leaderboardId string) (*LeaderboardResponse, error) {
	if leaderboardId == "" {
		return nil, ErrLeaderboardIdEmpty
	}
	comp, found := storage.Competitions[leaderboardId]
	if !found {
		return nil, ErrCompetetionNotFound
	}
	return asLeaderboardResponse(comp), nil
}

var GetLeaderboardForPlayer = func(playerId string) (*LeaderboardResponse, error) {
	comp, err := getCompetition(playerId)
	if err != nil {
		return nil, err
	}
	if comp.StartedAt().IsZero() {
		return nil, nil
	} else {
		return asLeaderboardResponse(comp), nil
	}
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

func asLeaderboardResponse(comp model.ICompetition) *LeaderboardResponse {
	if comp == nil {
		return nil
	}

	leaderboard := make([]PlayerScore, 0, len(comp.PlayersMap()))
	for _, player := range comp.Leaderboard() {
		leaderboard = append(leaderboard, PlayerScore{
			PlayerId: player.Player().Id(),
			Score:    player.Score(),
		})
	}

	return &LeaderboardResponse{
		Id:          comp.Id(),
		EndsAt:      comp.EndsAt(),
		Leaderboard: leaderboard,
	}
}

type LeaderboardResponse struct {
	Id          string        `json:"leaderboard_id"`
	EndsAt      time.Time     `json:"ends_at"`
	Leaderboard []PlayerScore `json:"leaderboard"`
}

type PlayerScore struct {
	PlayerId string `json:"player_id"`
	Score    int    `json:"score"`
}
