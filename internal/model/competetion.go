package model

import (
	"errors"
	"leaderboard/internal/config"
	"leaderboard/internal/timeprovider"
	"time"

	"github.com/google/uuid"
)

// ICompetition defines the contract for a competition
type ICompetition interface {
	Id() string
	CreatedAt() time.Time
	StartedAt() time.Time
	EndsAt() time.Time
	Players() []CompetingPlayer
	AddPlayer(player *Player) error
	Start() error
}

type Competition struct {
	id        string
	createdAt time.Time
	startedAt time.Time
	endsAt    time.Time
	players   []CompetingPlayer
}

var ErrCompetitionFull = errors.New("competition is full, cannot add more players")
var ErrCompetitionStarted = errors.New("competition has already started, cannot add players")
var ErrNotEnoughPlayers = errors.New("competition don't have enough players to start")

func NewCompetition() ICompetition {
	var comp = &Competition{
		id:        uuid.New().String(),
		createdAt: timeprovider.Current.Now(),
		startedAt: time.Time{},
		endsAt:    time.Time{},
		players:   make([]CompetingPlayer, 0, config.MaxPlayersForCompetition),
	}
	return comp
}

func (c *Competition) AddPlayer(player *Player) error {
	if len(c.players) >= config.MaxPlayersForCompetition {
		return ErrCompetitionFull
	}
	if !c.startedAt.IsZero() {
		return ErrCompetitionStarted
	}
	c.players = append(c.players, *NewCompetingPlayer(player))
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
	c.startedAt = timeprovider.Current.Now()
	c.endsAt = c.startedAt.Add(config.CompetitionDuration)
	return nil
}

func (c *Competition) Id() string {
	return c.id
}

func (c *Competition) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Competition) StartedAt() time.Time {
	return c.startedAt
}
func (c *Competition) EndsAt() time.Time {
	return c.endsAt
}
func (c *Competition) Players() []CompetingPlayer {
	return c.players
}
