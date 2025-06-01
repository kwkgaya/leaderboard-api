package matchmaking

import (
	"errors"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
	"sync"
	"time"
)

var (
	// TODO: Read from config or env
	MatchWaitDuration  = 30 * time.Second
	MatchRetryInterval = 1 * time.Second
)

var (
	ErrPlayerIdEmpty              = errors.New("player ID cannot be empty")
	ErrPlayerNotFound             = errors.New("player not found")
	ErrPlayerAlreadyInCompetition = errors.New("player is already in a competition")
)
var (
	// This mutex synchronizes the access to the waiting players and competitions maps
	// Also start competetion is accessed only by one goroutine at a time using this
	mutex = &sync.Mutex{}

	// Maps to hold players and competitions waiting for a match
	waitingPlayers      = make(map[uint]*model.Player)
	waitingCompetitions = make(map[uint]*model.Competition)
)

// JoinCompetition allows a player to join a competition.
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

	if player.ActiveCompetition() != nil {
		return nil, ErrPlayerAlreadyInCompetition
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
			comp, err := createNewCompetetion(player, waitingPlayer)
			if err != nil {
				return nil, err
			}
			return comp, nil
		} else {
			// No competition found, add player to waiting list
			waitingPlayers[player.Level()] = player
			// Start a timer to try starting a competition after the wait duration
			go func() {
				timer := time.NewTimer(MatchWaitDuration)
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

// TODO: This logic currently supports MinPlayersForCompetetion = 2 only
// Needs some updates to support higher values of MinPlayersForCompetetion
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
		higherLevel := player.Level() + uint(i)
		lowerLevel := player.Level() - uint(i)
		var waitingPlayer *model.Player
		if higherLevel <= model.MaxLevel {
			waitingPlayer, playerFound = waitingPlayers[higherLevel]
		}
		if !playerFound && lowerLevel >= model.MinLevel {
			waitingPlayer, playerFound = waitingPlayers[lowerLevel]
		}
		if playerFound {
			comp, err := createNewCompetetion(player, waitingPlayer)
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
		if higherLevel >= model.MaxLevel && lowerLevel <= model.MinLevel {
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
	ticker := time.NewTicker(MatchRetryInterval)
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

func createNewCompetetion(player *model.Player, waitingPlayer *model.Player) (*model.Competition, error) {
	if player == nil || waitingPlayer == nil {
		return nil, errors.New("both players must be provided")
	}

	comp := model.NewCompetition(player.Level())
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
