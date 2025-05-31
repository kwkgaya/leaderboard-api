package matchmaking

import (
	"errors"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
)

var ErrPlayerIdEmpty = errors.New("player ID cannot be empty")
var ErrPlayerNotFound = errors.New("player not found")
var ErrPlayerAlreadyInCompetition = errors.New("player is already in a competition")

var waitingPlayers = make(map[uint]*model.Player)
var waitingCompetitions = make(map[uint]*model.Competition)

func JoinCompetition(playerID string) (*model.Competition, error) {
	if playerID == "" {
		return nil, ErrPlayerIdEmpty
	}

	player := storage.Players[playerID]
	if player == nil {
		return nil, ErrPlayerNotFound
	}
	if player.ActiveCompetition() != nil {
		return nil, ErrPlayerAlreadyInCompetition
	}

	comp, compFound := waitingCompetitions[player.Level()]
	if compFound {
		comp.AddPlayer(player)
		return comp, nil
	} else {
		waitingPlayer, waitingPlayerFound := waitingPlayers[player.Level()]
		if waitingPlayerFound {
			// If there's already a waiting player at this level, create a new competition
			comp = model.NewCompetition(player.Level())
			storage.Competitions[comp.Id()] = comp
			err := comp.AddPlayer(player)
			if err != nil {
				return nil, err
			}
			err = comp.AddPlayer(waitingPlayer)
			if err != nil {
				return nil, err
			}
			waitingCompetitions[player.Level()] = comp
			delete(waitingPlayers, player.Level())
			return comp, nil
		} else {
			// No competition found, add player to waiting list
			waitingPlayers[player.Level()] = player
			return nil, nil // Player is now waiting for a match
		}
	}
}
