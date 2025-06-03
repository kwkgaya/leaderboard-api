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
	waitingCompetitions = make(map[int]model.ICompetition)
	// Slice to hold the competitions in the order they are created
	orderedCompetitions = make([]model.ICompetition, 0, config.MaxCompetitionsInMemory)
)

var JoinCompetition = func(playerID string) (model.ICompetition, error) {
	if playerID == "" {
		return nil, ErrPlayerIdEmpty
	}
	mutex.Lock()
	defer mutex.Unlock()

	player, playerFound := storage.Players[playerID]
	if !playerFound {
		return nil, ErrPlayerNotFound
	}
	comp := player.Competition()
	if comp != nil {
		// TODO: This is domain logic. Move to the model
		if !comp.StartedAt().IsZero() && comp.EndsAt().Before(timeprovider.Current.Now()) {
			// If the competition has ended, reset the player's competition
			player.SetCompetition(nil)
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
		comp, err := createNewCompetition(player)
		if err != nil {
			return nil, err
		}
		// Start a timer to try starting a competition after the wait duration
		go func() {
			timer := time.NewTimer(config.MatchWaitDuration)
			<-timer.C
			err1 := tryStartCompetition(player)
			if err1 != nil {
				panic(err1) // Handle error appropriately in production code
			}
		}()

		return comp, nil // Player is now waiting for a match
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

	comp := player.Competition()
	if comp != nil && len(comp.PlayersMap()) >= config.MinPlayersForCompetition {
		// Player is already in a competition. Start it if not already started
		if comp.StartedAt().IsZero() {
			err := comp.Start()
			delete(waitingCompetitions, player.Level())
			if err != nil {
				return err
			}
		}
		return nil
	}

	matched := false
	for i := 1; ; i++ {
		// Try finding a matching competioion at closest levels
		higherLevel := player.Level() + i
		lowerLevel := player.Level() - i

		// Check if we have a competition waiting for a player at the higher or lower level
		var waitingComp model.ICompetition
		if higherLevel <= config.MaxLevel {
			if waitingComp = waitingCompetitions[higherLevel]; waitingComp != nil {
				err := waitingComp.AddPlayer(player)
				if err != nil {
					return err
				}
				matched = true
			}
		}
		if !matched && lowerLevel >= config.MinLevel {
			if waitingComp = waitingCompetitions[lowerLevel]; waitingComp != nil {
				err := waitingComp.AddPlayer(player)
				if err != nil {
					return err
				}
				delete(waitingCompetitions, lowerLevel)
				matched = true
			}
		}
		if waitingComp != nil {
			comp = waitingComp
			player.SetCompetition(comp)
			delete(waitingCompetitions, comp.InitialLevel())
		}

		if matched {
			err := comp.Start()

			if err != nil {
				return err
			}
			delete(waitingCompetitions, player.Level())
			break
		}
		if higherLevel >= config.MaxLevel && lowerLevel <= config.MinLevel {
			break
		}
	}

	// If still no matching player is found, we can start a ticker to keep checking
	if !matched {
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
		if player.Competition() != nil {
			ticker.Stop()
			return
		}
	}
}

func createNewCompetition(player *model.Player) (model.ICompetition, error) {
	if player == nil {
		return nil, errors.New("player must be provided")
	}

	comp := model.NewCompetition(player.Level())
	storage.Competitions[comp.Id()] = comp
	err := comp.AddPlayer(player)
	if err != nil {
		return nil, err
	}
	waitingCompetitions[player.Level()] = comp

	// TODO: Check if operations on this slice are performance optimum
	orderedCompetitions = append(orderedCompetitions, comp)

	ensureMaxCompetitionsInMemory()

	return comp, nil
}

func ensureMaxCompetitionsInMemory() {
	index := 0
	// Usually this will only delete the oldest competition that has started and ended
	// But ensure that we don't delete competitions that are still ongoing
	for len(storage.Competitions) > config.MaxCompetitionsInMemory &&
		!orderedCompetitions[index].StartedAt().IsZero() &&
		orderedCompetitions[index].EndsAt().Before(timeprovider.Current.Now()) {
		// Remove the oldest competition that has started and ended
		delete(storage.Competitions, orderedCompetitions[index].Id())
		index += 1
	}
	if index > 0 {
		orderedCompetitions = orderedCompetitions[index:]
	}
}
