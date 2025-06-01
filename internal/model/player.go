package model

import (
	"fmt"
	"leaderboard/internal/config"
)

type Player struct {
	id                string
	level             int
	countryCode       string
	activeCompetition ICompetition
}

func NewPlayer(id string, level int, countryCode string) *Player {
	if level < config.MinLevel || level > config.MaxLevel {
		panic("player level must be between MinLevel " + fmt.Sprint(config.MinLevel) + " and MaxLevel " + fmt.Sprint(config.MaxLevel))
	}
	return &Player{
		id:          id,
		level:       level,
		countryCode: countryCode,
	}
}
func (p *Player) Id() string {
	return p.id
}
func (p *Player) Level() int {
	return p.level
}
func (p *Player) CountryCode() string {
	return p.countryCode
}
func (p *Player) Competition() ICompetition {
	return p.activeCompetition
}
func (p *Player) SetCompetition(c ICompetition) {
	p.activeCompetition = c
}
