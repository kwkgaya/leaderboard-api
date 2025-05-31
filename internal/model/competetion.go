package model

import (
	"errors"
	"leaderboard/internal/timeprovider"
	"time"

	"github.com/google/uuid"
)

const (
	MaxPlayersForCompetetion int = 10
	MinPlayersForCompetetion int = 2
	CompetitionDuration          = 1 * time.Hour
)

type Competition struct {
	id           string
	initialLevel uint
	createdAt    time.Time
	startedAt    time.Time
	endsAt       time.Time
	players      []CompetingPlayer
}

var ErrCompetitionFull = errors.New("competition is full, cannot add more players")
var ErrCompetitionStarted = errors.New("competition has already started, cannot add players")
var ErrNotEnoughPlayers = errors.New("competetion has less than two players")

func NewCompetition(initialLevel uint) *Competition {
	var comp = &Competition{
		id:           uuid.New().String(),
		initialLevel: initialLevel,
		createdAt:    timeprovider.Current.Now(),
		startedAt:    time.Time{},
		endsAt:       time.Time{},
		players:      make([]CompetingPlayer, 0, MaxPlayersForCompetetion),
	}
	return comp
}

func (c *Competition) AddPlayer(player *Player) error {
	if len(c.players) >= MaxPlayersForCompetetion {
		return ErrCompetitionFull
	}
	if !c.startedAt.IsZero() {
		return ErrCompetitionStarted
	}
	c.players = append(c.players, *NewCompetingPlayer(player))
	player.SetActiveCompetition(c)

	if len(c.players) == MaxPlayersForCompetetion {
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
	if len(c.players) < MinPlayersForCompetetion {
		return ErrNotEnoughPlayers
	}
	c.startedAt = timeprovider.Current.Now()
	c.endsAt = c.startedAt.Add(CompetitionDuration)
	return nil
}

func (c *Competition) Id() string {
	return c.id
}

func (c *Competition) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Competition) InitialLevel() uint {
	return c.initialLevel
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
