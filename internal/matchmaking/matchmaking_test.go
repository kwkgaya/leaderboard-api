package matchmaking

import (
	"errors"
	"leaderboard/internal/storage"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	storage.AddPlayers([]storage.NewPlayer{
		{Id: "alice_smith", CountryCode: "US", Level: 1},
		{Id: "bob_jones", CountryCode: "GB", Level: 2},
		{Id: "carlos_mendez", CountryCode: "MX", Level: 3},
	})
}

func TestJoinCompetitionBasic(t *testing.T) {
	// Basic test cases for joining a competition
	tests := []struct {
		name          string
		playerId      string
		expectedError error
	}{
		{"Empty player Id", "", ErrPlayerIdEmpty},
		{"Unknown player Id", "unknown", ErrPlayerNotFound},
		{"Valid player Id", "alice_smith", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := JoinCompetition(tt.playerId)
			if err == nil && tt.expectedError == nil {
				return
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("JoinCompetition() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
