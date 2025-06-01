package model

type CompetingPlayer struct {
	player *Player
	score  int
}

func NewCompetingPlayer(player *Player) *CompetingPlayer {
	return &CompetingPlayer{
		player: player,
		score:  0}
}

func (p *CompetingPlayer) Score() int {
	return p.score
}

func (p *CompetingPlayer) Player() *Player {
	return p.player
}

func (p *CompetingPlayer) AddScore(score int) {
	p.score += score
}
