package model

type Player struct {
	ID          string
	Level       int
	activeCompetition *Competition
}
