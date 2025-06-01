package model

import "leaderboard/internal/config"

type Player struct {
	id                string
	level             int
	countryCode       string
	activeCompetition *Competition
}

func NewPlayer(id string, level int, countryCode string) *Player {
	if level < config.MinLevel || level > config.MaxLevel {
		panic("player level must be between MinLevel and MaxLevel")
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
func (p *Player) ActiveCompetition() *Competition {
	return p.activeCompetition
}
func (p *Player) SetActiveCompetition(c *Competition) {
	p.activeCompetition = c
}
