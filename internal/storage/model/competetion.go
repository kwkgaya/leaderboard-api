package model

import (
	"time"
)

type Competition struct {
	ID        string
	CreatedAt time.Time
	EndsAt    time.Time
	Players   [10]*PlayerCompetition
}