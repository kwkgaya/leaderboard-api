package model

import (
	"leaderboard/internal/timeprovider"
	"time"
)

type CompetingPlayer struct {
	player    *Player
	score     uint
	createdAt time.Time
}

func NewCompetingPlayer(player *Player) *CompetingPlayer {
	return &CompetingPlayer{
		player:    player,
		createdAt: timeprovider.Current.Now(),
		score:     0}
}

func (p *CompetingPlayer) CreatedAt() time.Time {
	return p.createdAt
}

func (p *CompetingPlayer) Score() uint {
	return p.score
}

func (p *CompetingPlayer) Player() *Player {
	return p.player
}

func (p *CompetingPlayer) AddScore(score uint) {
	p.score += score
}
