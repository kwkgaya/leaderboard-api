package model

type Player struct {
	id                string
	level             uint
	countryCode       string
	activeCompetition *Competition
}

func NewPlayer(id string, level uint, countryCode string) *Player {
	return &Player{
		id:          id,
		level:       level,
		countryCode: countryCode,
	}
}
func (p *Player) Id() string {
	return p.id
}
func (p *Player) Level() uint {
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
