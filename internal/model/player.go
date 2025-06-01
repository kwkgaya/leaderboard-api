package model

const (
	MaxLevel = 10 // Maximum level a player can have
	MinLevel = 1  // Minimum level a player can have
)

type Player struct {
	id                string
	level             int
	countryCode       string
	activeCompetition *Competition
}

func NewPlayer(id string, level int, countryCode string) *Player {
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
