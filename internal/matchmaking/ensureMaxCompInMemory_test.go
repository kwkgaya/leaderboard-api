package matchmaking

import (
	"leaderboard/internal/config"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
	"testing"
	"time"
)

func tearDownEnsureMaxCompetitionsInMemory() {
	config.MaxCompetitionsInMemory = 100
	orderedCompetitions = make([]model.ICompetition, 0)

	clear(waitingCompetitions)
	clear(orderedCompetitions)
	clear(storage.Players)
	clear(storage.Competitions)
}

func TestEnsureMaxCompetitionsInMemory_RemovesOldEndedCompetitions(t *testing.T) {

	// Set config to allow only 2 competitions in memory
	config.MaxCompetitionsInMemory = 2

	// Create 4 competitions, 3 of which are ended
	now := time.Now()
	for i := 0; i < 4; i++ {
		comp := model.NewCompetition(1).(*model.Competition)
		storage.Competitions[comp.Id()] = comp
		orderedCompetitions = append(orderedCompetitions, comp)
		// Simulate started and ended competitions for first 3
		if i < 3 {
			comp.SetStartedAt(now.Add(-time.Duration(i+2) * time.Minute))
			comp.SetEndsAt(now.Add(-time.Duration(i+1) * time.Minute))
		}
	}

	ensureMaxCompetitionsInMemory()

	// Only the last 2 competitions should remain in orderedCompetitions
	if len(orderedCompetitions) != 2 {
		t.Errorf("expected 2 competitions in memory, got %d", len(orderedCompetitions))
	}
	// The competitions in storage should match the ones in orderedCompetitions
	for _, comp := range orderedCompetitions {
		if _, ok := storage.Competitions[comp.Id()]; !ok {
			t.Errorf("competition %s should remain in storage", comp.Id())
		}
	}

	tearDownEnsureMaxCompetitionsInMemory()
}

func TestEnsureMaxCompetitionsInMemory_DoesNotRemoveOngoingCompetitions(t *testing.T) {

	config.MaxCompetitionsInMemory = 1

	// Create 3 competitions, one ended, two ongoing
	compEnded := model.NewCompetition(1).(*model.Competition)
	compEnded.SetStartedAt(time.Now().Add(-2 * time.Minute))
	compEnded.SetEndsAt(time.Now().Add(-1 * time.Minute))
	storage.Competitions[compEnded.Id()] = compEnded

	compOngoing := model.NewCompetition(1).(*model.Competition)
	compEnded.SetStartedAt(time.Now().Add(-2 * time.Minute))
	compOngoing.SetEndsAt(time.Now().Add(10 * time.Minute))
	storage.Competitions[compOngoing.Id()] = compOngoing

	compOngoing2 := model.NewCompetition(1).(*model.Competition)
	compEnded.SetStartedAt(time.Now().Add(-2 * time.Minute))
	compOngoing2.SetEndsAt(time.Now().Add(5 * time.Minute))
	storage.Competitions[compOngoing2.Id()] = compOngoing2

	orderedCompetitions = append(orderedCompetitions, compEnded, compOngoing, compOngoing2)

	ensureMaxCompetitionsInMemory()

	// Ongoing competition should not be removed
	if _, ok := storage.Competitions[compOngoing.Id()]; !ok {
		t.Errorf("ongoing competition should not be removed from storage")
	}
	// Ended competition should be removed
	if _, ok := storage.Competitions[compEnded.Id()]; ok {
		t.Errorf("ended competition should be removed from storage")
	}
	// Only ongoing competition should remain in orderedCompetitions
	if len(orderedCompetitions) != 2 ||
		orderedCompetitions[0].Id() != compOngoing.Id() ||
		orderedCompetitions[1].Id() != compOngoing2.Id() {
		t.Errorf("only ongoing competitions should remain in orderedCompetitions")
	}

	tearDownEnsureMaxCompetitionsInMemory()
}

func TestEnsureMaxCompetitionsInMemory_NoRemovalIfBelowLimit(t *testing.T) {

	config.MaxCompetitionsInMemory = 5

	// Create 3 competitions, all ended
	for i := 0; i < 3; i++ {
		comp := model.NewCompetition(1).(*model.Competition)
		comp.SetEndsAt(time.Now().Add(-time.Duration(i+2) * time.Minute))
		comp.SetEndsAt(time.Now().Add(-time.Duration(i+1) * time.Minute))
		storage.Competitions[comp.Id()] = comp
		orderedCompetitions = append(orderedCompetitions, comp)
	}

	ensureMaxCompetitionsInMemory()

	// All competitions should remain
	if len(orderedCompetitions) != 3 {
		t.Errorf("expected 3 competitions in memory, got %d", len(orderedCompetitions))
	}
	for _, comp := range orderedCompetitions {
		if _, ok := storage.Competitions[comp.Id()]; !ok {
			t.Errorf("competition %s should remain in storage", comp.Id())
		}
	}

	tearDownEnsureMaxCompetitionsInMemory()
}

func TestEnsureMaxCompetitionsInMemory_HandlesNoCompetitions(t *testing.T) {

	config.MaxCompetitionsInMemory = 2

	orderedCompetitions = nil
	ensureMaxCompetitionsInMemory()

	if len(orderedCompetitions) != 0 {
		t.Errorf("expected 0 competitions, got %d", len(orderedCompetitions))
	}

	tearDownEnsureMaxCompetitionsInMemory()
}
