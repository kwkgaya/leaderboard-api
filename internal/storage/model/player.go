package model

type Player struct {
	ID                string
	Level             int
	CountryCode       string
	ActiveCompetition *Competition
}
