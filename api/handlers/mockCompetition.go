package handlers

import (
	"leaderboard/internal/config"
	"leaderboard/internal/model"
	"leaderboard/internal/timeprovider"
	"time"
)

// mockCompetition implements the minimal interface needed for JoinHandler
type mockCompetition struct {
	id        string
	startedAt time.Time
	endsAt    time.Time
}

func (m *mockCompetition) Id() string           { return m.id }
func (m *mockCompetition) StartedAt() time.Time { return m.startedAt }
func (m *mockCompetition) EndsAt() time.Time    { return m.endsAt }
func (m *mockCompetition) CreatedAt() time.Time { return time.Time{} }
func (m *mockCompetition) Players() []model.CompetingPlayer {
	return nil // Not needed for these tests
}
func (m *mockCompetition) AddPlayer(player *model.Player) error {
	return nil // Not needed for these tests
}
func (m *mockCompetition) Start() error {
	m.startedAt = timeprovider.Current.Now()
	m.endsAt = m.startedAt.Add(config.CompetitionDuration)
	return nil

}
