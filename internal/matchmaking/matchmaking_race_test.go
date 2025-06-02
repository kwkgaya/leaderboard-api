package matchmaking

import (
	"leaderboard/internal/storage"
	"sync"
	"testing"
)

func TestJoinCompetition_RaceCondition(t *testing.T) {
	var wg sync.WaitGroup
	playerCount := 1000

	// Use a map to check for duplicate joins
	joined := make(map[string]bool)
	var mu sync.Mutex

	// AddPlayers is not thread-safe, Thefore add players before starting the goroutines
	for i := 0; i < playerCount; i++ {
		playerID := "player" + string(rune(i))
		storage.AddPlayers([]storage.NewPlayer{
			{Id: playerID, CountryCode: "US", Level: i%3 + 1},
		})
	}

	for i := 0; i < playerCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			playerID := "player" + string(rune(i))

			if _, err := JoinCompetition(playerID); err != nil {
				t.Errorf("unexpected error for %s: %v", playerID, err)
			}
			mu.Lock()
			if joined[playerID] {
				t.Errorf("player %s joined more than once", playerID)
			}
			joined[playerID] = true
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	// Check that all players joined
	if len(joined) != playerCount {
		t.Errorf("expected %d players joined, got %d", playerCount, len(joined))
	}
}
