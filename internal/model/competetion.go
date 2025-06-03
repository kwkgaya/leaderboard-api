package model

import (
	"errors"
	"leaderboard/internal/config"
	"leaderboard/internal/timeprovider"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ICompetition defines the contract for a competition
type ICompetition interface {
	Id() string
	StartedAt() time.Time
	EndsAt() time.Time
	PlayersMap() map[string]*CompetingPlayer
	Leaderboard() []*CompetingPlayer
	AddPlayer(player *Player) error
	Start() error
	AddScore(playerId string, points int) error
	InitialLevel() int
}

type Competition struct {
	id            string
	startedAt     time.Time
	endsAt        time.Time
	players       map[string]*CompetingPlayer
	sortedPlayers []*CompetingPlayer
	scoreMutex    *sync.Mutex
	initialLevel  int
}

var (
	ErrCompetitionFull            = errors.New("competition is full, cannot add more players")
	ErrCompetitionStarted         = errors.New("competition has already started, cannot add players")
	ErrCompetitionNotStarted      = errors.New("competition has not started yet")
	ErrNotEnoughPlayers           = errors.New("competition don't have enough players to start")
	ErrPlayerAlreadyInCompetition = errors.New("player is already in this competition")

	ErrPlayerIdEmpty  = errors.New("player ID cannot be empty")
	ErrPlayerNotFound = errors.New("player not found in competition")
	ErrPointsNegative = errors.New("points cannot be negative")
)

var (
	competetionsStarted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "leaderboard_competitions_started_total",
		Help: "The total number of competitions started",
	})
)

func NewCompetition(initialLevel int) ICompetition {
	var comp = &Competition{
		id:           uuid.New().String(),
		initialLevel: initialLevel,
		startedAt:    time.Time{},
		endsAt:       time.Time{},
		players:      make(map[string]*CompetingPlayer, config.MaxPlayersForCompetition),
	}
	return comp
}

func (c *Competition) AddPlayer(player *Player) error {
	if player == nil {
		return ErrPlayerIdEmpty
	}
	if len(c.players) >= config.MaxPlayersForCompetition {
		return ErrCompetitionFull
	}
	if !c.startedAt.IsZero() {
		return ErrCompetitionStarted
	}
	if c.players[player.Id()] != nil {
		return ErrPlayerAlreadyInCompetition
	}
	c.players[player.Id()] = &CompetingPlayer{
		player: player,
		score:  0,
	}
	player.SetCompetition(c)

	if len(c.players) == config.MaxPlayersForCompetition {
		if err := c.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Competition) Start() error {
	if !c.startedAt.IsZero() {
		return ErrCompetitionStarted
	}
	if len(c.players) < config.MinPlayersForCompetition {
		return ErrNotEnoughPlayers
	}
	c.sortedPlayers = slices.Collect(maps.Values(c.players))
	c.scoreMutex = &sync.Mutex{}

	c.startedAt = timeprovider.Current.Now()
	c.endsAt = c.startedAt.Add(config.CompetitionDuration)
	competetionsStarted.Inc()
	return nil
}

func (c *Competition) AddScore(playerId string, points int) error {
	if playerId == "" {
		return ErrPlayerIdEmpty
	}
	if points < 0 {
		return ErrPointsNegative
	}
	if c.startedAt.IsZero() {
		return ErrCompetitionNotStarted
	}

	if compPlayer, found := c.players[playerId]; found {
		c.scoreMutex.Lock()
		defer c.scoreMutex.Unlock()

		compPlayer.AddScore(points)
		slices.SortStableFunc(c.sortedPlayers, func(a, b *CompetingPlayer) int {
			if a.Score() == b.Score() {
				return strings.Compare(a.Player().Id(), b.Player().Id())
			} else {
				return b.Score() - a.Score()
			}
		})
		return nil
	} else {
		return ErrPlayerNotFound
	}
}

func (c *Competition) Id() string {
	return c.id
}

func (c *Competition) StartedAt() time.Time {
	return c.startedAt
}
func (c *Competition) EndsAt() time.Time {
	return c.endsAt
}
func (c *Competition) PlayersMap() map[string]*CompetingPlayer {
	return c.players
}
func (c *Competition) Leaderboard() []*CompetingPlayer {
	return c.sortedPlayers
}
func (c *Competition) InitialLevel() int {
	return c.initialLevel
}

func (c *Competition) SetStartedAt(time time.Time) {
	c.startedAt = time
}

func (c *Competition) SetEndsAt(time time.Time) {
	c.endsAt = time
}
