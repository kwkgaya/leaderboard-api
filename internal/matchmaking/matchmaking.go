package matchmaking

import (
	"errors"
	"leaderboard/internal/config"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
	"leaderboard/internal/timeprovider"
	"sync"
	"time"
)

var (
	ErrPlayerIdEmpty              = errors.New("player ID cannot be empty")
	ErrPlayerNotFound             = errors.New("player not found")
	ErrPlayerAlreadyInCompetition = errors.New("player is already in a competition")
)
var (
	// This mutex synchronizes the access to the waiting players and competitions maps
	// Also start competition is accessed only by one goroutine at a time using this
	mutex = &sync.Mutex{}

	// Maps to hold players and competitions waiting for a match
	waitingPlayers      = make(map[int]*model.Player)
	waitingCompetitions = make(map[int]*model.Competition)
)

func JoinCompetition(playerID string) (*model.Competition, error) {
	if playerID == "" {
		return nil, ErrPlayerIdEmpty
	}
	mutex.Lock()
	defer mutex.Unlock()

	player := storage.Players[playerID]
	if player == nil {
		return nil, ErrPlayerNotFound
	}
	activeComp := player.ActiveCompetition()
	if activeComp != nil {
		// TODO: This is domain logic. Move to the model
		if !activeComp.StartedAt().IsZero() && activeComp.EndsAt().Before(timeprovider.Current.Now()) {
			// If the active competition has ended, reset the player's active competition
			player.SetActiveCompetition(nil)
		} else {
			return nil, ErrPlayerAlreadyInCompetition
		}
	}

	comp, compFound := waitingCompetitions[player.Level()]
	if compFound {
		comp.AddPlayer(player)

		// Competition may start immediately if it has enough players
		if !comp.StartedAt().IsZero() {
			delete(waitingCompetitions, player.Level())
		}
		return comp, nil
	} else {
		waitingPlayer, waitingPlayerFound := waitingPlayers[player.Level()]
		if waitingPlayerFound {
			// If there's already a waiting player at this level, create a new competition and add both players
			comp, err := createNewCompetition(player, waitingPlayer)
			if err != nil {
				return nil, err
			}
			return comp, nil
		} else {
			// No competition found, add player to waiting list
			waitingPlayers[player.Level()] = player
			// Start a timer to try starting a competition after the wait duration
			go func() {
				timer := time.NewTimer(config.MatchWaitDuration)
				<-timer.C
				err1 := tryStartCompetition(player)
				if err1 != nil {
					panic(err1) // Handle error appropriately in production code
				}
			}()

			return nil, nil // Player is now waiting for a match
		}
	}
}

// TODO: This logic currently supports MinPlayersForCompetition = 2 only
// Needs some updates to support higher values of MinPlayersForCompetition
func tryStartCompetition(player *model.Player) error {
	if player == nil {
		panic("player cannot be nil")
	}

	mutex.Lock()
	defer mutex.Unlock()

	if player.ActiveCompetition() != nil {
		// Player is already in a competition. Start it if not already started
		activeComp := player.ActiveCompetition()
		if activeComp.StartedAt().IsZero() {
			err := activeComp.Start()
			delete(waitingCompetitions, player.Level())
			if err != nil {
				return err
			}
		}
		return nil
	}

	playerFound := false
	for i := 1; ; i++ {
		// Try finding a matching player at closest levels
		higherLevel := player.Level() + i
		lowerLevel := player.Level() - i
		var waitingPlayer *model.Player
		if higherLevel <= config.MaxLevel {
			waitingPlayer, playerFound = waitingPlayers[higherLevel]
		}
		if !playerFound && lowerLevel >= config.MinLevel {
			waitingPlayer, playerFound = waitingPlayers[lowerLevel]
		}
		if playerFound {
			comp, err := createNewCompetition(player, waitingPlayer)
			if err != nil {
				return err
			}

			err = comp.Start()
			delete(waitingCompetitions, player.Level())
			delete(waitingCompetitions, waitingPlayer.Level())

			if err != nil {
				return err
			}

			break
		}
		if higherLevel >= config.MaxLevel && lowerLevel <= config.MinLevel {
			break
		}
	}

	// If still no matching player is found, we can start a ticker to keep checking
	if !playerFound {
		go scheduleTickerForPlayer(player)
	}

	// If we reach here, it means no competition was started but player is still in the waiting list
	// Ticker will keep trying to find a match
	return nil
}

func scheduleTickerForPlayer(player *model.Player) {
	ticker := time.NewTicker(config.MatchRetryInterval)
	// TODO: Stop retry after a certain number of attempts or time limit
	for range ticker.C {
		err := tryStartCompetition(player)
		if err != nil {
			panic(err) // Handle error appropriately in production code
		}
		if player.ActiveCompetition() != nil {
			ticker.Stop()
			return
		}
	}
}

func createNewCompetition(player *model.Player, waitingPlayer *model.Player) (*model.Competition, error) {
	if player == nil || waitingPlayer == nil {
		return nil, errors.New("both players must be provided")
	}

	comp := model.NewCompetition()
	storage.Competitions[comp.Id()] = comp
	err := comp.AddPlayer(player)
	if err != nil {
		return nil, err
	}
	err = comp.AddPlayer(waitingPlayer)
	if err != nil {
		return nil, err
	}
	waitingCompetitions[waitingPlayer.Level()] = comp

	delete(waitingPlayers, waitingPlayer.Level())
	delete(waitingPlayers, player.Level())

	return comp, nil
}
