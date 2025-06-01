package matchmaking

import (
	"errors"
	"fmt"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
	"testing"
	"time"
)

func setup() {
	// Reset global state before test
	clear(waitingPlayers)
	clear(waitingCompetitions)
	storage.Players = map[string]*model.Player{}
	storage.Competitions = map[string]*model.Competition{}

	MatchWaitDuration = 30 * time.Second // Set a short wait duration for testing

	storage.AddPlayers([]storage.NewPlayer{
		{Id: "alice", CountryCode: "US", Level: 1},
		{Id: "bob", CountryCode: "GB", Level: 2},
		{Id: "carlos", CountryCode: "MX", Level: 3},
		{Id: "alice_1", CountryCode: "IN", Level: 1},
		{Id: "bob_1", CountryCode: "IN", Level: 2},
		{Id: "carlos_1", CountryCode: "IN", Level: 3},
		{Id: "alice_2", CountryCode: "IN", Level: 1},
		{Id: "bob_2", CountryCode: "IN", Level: 2},
		{Id: "carlos_2", CountryCode: "IN", Level: 3},
		{Id: "ian", CountryCode: "IN", Level: 10},
	})
}

func TestJoinCompetitionBasic(t *testing.T) {
	setup()
	// Basic test cases for joining a competition
	tests := []struct {
		name          string
		playerId      string
		expectedError error
	}{
		{"Empty player Id", "", ErrPlayerIdEmpty},
		{"Unknown player Id", "unknown", ErrPlayerNotFound},
		{"Valid player Id", "alice", nil},
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

func TestJoinCompetition_Matchmaking(t *testing.T) {
	setup()

	// First player joins, should be put in waiting list
	comp, err := JoinCompetition("bob")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp != nil {
		t.Errorf("expected nil competition for first player, got %v", comp)
	}
	if waitingPlayers[2] == nil || waitingPlayers[2].Id() != "bob" {
		t.Errorf("bob should be in waitingPlayers at level 2")
	}

	// Second player joins, should create a competition
	comp, err = JoinCompetition("bob_1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp == nil {
		t.Fatalf("expected competition to be created for bob_1")
	}
	if len(comp.Players()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.Players()))
	}
	if waitingPlayers[2] != nil {
		t.Errorf("waitingPlayers at level 2 should be empty after match")
	}
	if !comp.StartedAt().IsZero() {
		t.Errorf("competition should not have started yet, got started at %v", comp.StartedAt())
	}
	if len(waitingCompetitions) != 1 {
		t.Errorf("waitingCompetitions should have 1 competition, got %d", len(waitingCompetitions))
	}
}

func TestJoinCompetition_AlreadyInCompetition(t *testing.T) {
	setup()

	storage.AddPlayers([]storage.NewPlayer{
		{Id: "player3", CountryCode: "US", Level: 7},
	})

	player := storage.Players["player3"]
	// Simulate player already in a competition
	fakeComp := model.NewCompetition(player.Level())
	player.SetActiveCompetition(fakeComp)

	comp, err := JoinCompetition("player3")

	if !errors.Is(err, ErrPlayerAlreadyInCompetition) {
		t.Errorf("expected ErrPlayerAlreadyInCompetition, got %v", err)
	}
	if comp != nil {
		t.Errorf("expected nil competition, got %v", comp)
	}
}

func TestJoinCompetition_JoinMaxplayers_CompetetionStarted(t *testing.T) {
	setup()
	var previousComp *model.Competition
	for i := 1; i <= model.MaxPlayersForCompetetion; i++ {
		playerId := fmt.Sprintf("player%v", i)
		storage.AddPlayers([]storage.NewPlayer{
			{Id: playerId, CountryCode: "US", Level: 5},
		})

		comp, err := JoinCompetition(playerId)

		if i == 1 {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if comp != nil {
				t.Fatalf("expected competition to be not created for %s", playerId)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if comp == nil {
				t.Fatalf("expected competition to be returned for %s", playerId)
			}
			if len(comp.Players()) != i {
				t.Errorf("competition should have %d players, got %d", i, len(comp.Players()))
			}
			if previousComp != nil && comp.Id() != previousComp.Id() {
				t.Errorf("expected same competition for player %s, got different competition %s", playerId, comp.Id())
			}
			previousComp = comp

			player := storage.Players[playerId]
			if player.ActiveCompetition() == nil || player.ActiveCompetition().Id() != comp.Id() {
				t.Errorf("player %s should be in competition %s, got %v", playerId, comp.Id(), player.ActiveCompetition())
			}

			if i == model.MaxPlayersForCompetetion {
				if comp.StartedAt().IsZero() {
					t.Errorf("competition should have started after adding %d players, got started at %v", model.MaxPlayersForCompetetion, comp.StartedAt())
				}
				if len(waitingCompetitions) != 0 {
					t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
				}

				// Test player1 only one time
				player1Id := fmt.Sprintf("player%v", 1)
				player1 := storage.Players[player1Id]
				if player1.ActiveCompetition() == nil || player1.ActiveCompetition().Id() != comp.Id() {
					t.Errorf("player %s should be in competition %s, got %v", player1Id, comp.Id(), player.ActiveCompetition())
				}
			} else {
				if !comp.StartedAt().IsZero() {
					t.Errorf("competition should not have started yet, got started at %v", comp.StartedAt())
				}
				if len(waitingCompetitions) != 1 {
					t.Errorf("waitingCompetitions should have 1 competition, got %d", len(waitingCompetitions))
				}
			}
		}
	}
}

func TestJoinCompetition_JoinMaxplayersAndTwoMore_NewCompetetionStarted(t *testing.T) {
	setup()
	var previousComp *model.Competition
	for i := 1; i <= model.MaxPlayersForCompetetion+2; i++ {
		playerId := fmt.Sprintf("player%v", i)
		storage.AddPlayers([]storage.NewPlayer{
			{Id: playerId, CountryCode: "US", Level: 5},
		})

		comp, err := JoinCompetition(playerId)

		if i == model.MaxPlayersForCompetetion {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			previousComp = comp
		} else if i == model.MaxPlayersForCompetetion+1 {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if comp != nil {
				t.Fatalf("expected new competition to be not created for %s", playerId)
			}
		} else if i == model.MaxPlayersForCompetetion+2 {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if comp.Id() == previousComp.Id() {
				t.Fatalf("expected new competition to be created for %s, got same competition %s", playerId, comp.Id())
			}
		}
	}
}

func TestJoinCompetition_MatchedwithTwoPlayers_CompetetionStartsAfterMatchWaitDuration(t *testing.T) {
	setup()
	MatchWaitDuration = 500 * time.Millisecond // Set a short wait duration for testing

	_, err := JoinCompetition("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	comp, err := JoinCompetition("bob_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !comp.StartedAt().IsZero() {
		t.Errorf("competition should not have started immediately, got started at %v", comp.StartedAt())
	}

	time.Sleep(1 * time.Second) // Wait for starting competition after MatchWaitDuration
	if comp.StartedAt().IsZero() {
		t.Errorf("competition should have started after %v, got started at %v", MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.Players()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.Players()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
}

func TestJoinCompetition_MatchedwithTwoPlayersInTwoLevels_CompetetionStartsAfterMatchWaitDuration(t *testing.T) {
	setup()
	MatchWaitDuration = 1 * time.Second // Set a short wait duration for testing

	_, err := JoinCompetition("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = JoinCompetition("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alice := storage.Players["alice"]
	bob := storage.Players["bob"]

	time.Sleep(2 * time.Second) // Wait for starting competition after MatchWaitDuration

	if alice.ActiveCompetition() == nil {
		t.Errorf("alice should be in a competition, got nil")
	}
	if bob.ActiveCompetition() == nil {
		t.Errorf("bob should be in a competition, got nil")
	}
	if alice.ActiveCompetition().Id() != bob.ActiveCompetition().Id() {
		t.Errorf("alice and bob should be in the same competition, got %s and %s", alice.ActiveCompetition().Id(), bob.ActiveCompetition().Id())
	}
	comp := alice.ActiveCompetition()

	if comp.StartedAt().IsZero() {
		t.Errorf("competition should have started after %v, got started at %v", MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.Players()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.Players()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
}

func TestJoinCompetition_MatchedwithTwoPlayersInMinandMaxLevels_CompetetionStartsAfterMatchWaitDuration(t *testing.T) {
	setup()
	MatchWaitDuration = 1 * time.Second // Set a short wait duration for testing

	_, err := JoinCompetition("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = JoinCompetition("ian")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alice := storage.Players["alice"]
	ian := storage.Players["ian"]

	time.Sleep(2 * time.Second) // Wait for starting competition after MatchWaitDuration

	if alice.ActiveCompetition() == nil {
		t.Errorf("alice should be in a competition, got nil")
	}
	if ian.ActiveCompetition() == nil {
		t.Errorf("ian should be in a competition, got nil")
	}
	if alice.ActiveCompetition().Id() != ian.ActiveCompetition().Id() {
		t.Errorf("alice and bob should be in the same competition, got %s and %s", alice.ActiveCompetition().Id(), ian.ActiveCompetition().Id())
	}
	comp := alice.ActiveCompetition()

	if comp.StartedAt().IsZero() {
		t.Errorf("competition should have started after %v, got started at %v", MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.Players()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.Players()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
}

func TestJoinCompetition_CompetetionStartAfterWait_NewJoineesAddedToNewCompetetion(t *testing.T) {
	setup()
	MatchWaitDuration = 1500 * time.Millisecond // Set a short wait duration for testing

	_, err := JoinCompetition("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = JoinCompetition("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alice := storage.Players["alice"]

	time.Sleep(2 * time.Second) // Wait for starting competition after MatchWaitDuration
	comp1 := alice.ActiveCompetition()

	alice1comp, err := JoinCompetition("alice_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alice1comp != nil {
		t.Errorf("expected alice_1 to be waiting for a match, got competition %s", alice1comp.Id())
	}

	bob1comp, err := JoinCompetition("bob_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bob1comp != nil {
		t.Errorf("expected bob_1 to be waiting for a match, got competition %s", bob1comp.Id())
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should have 0 competitions, got %d", len(waitingCompetitions))
	}

	alice2comp, err := JoinCompetition("alice_2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alice2comp == nil {
		t.Fatalf("expected alice_2 to be added to competition, got nil")
	}
	if alice2comp.Id() == comp1.Id() {
		t.Errorf("expected alice_2 to be added to a new competition, but added to old one %s", alice2comp.Id())
	}
}

func TestJoinCompetition_NotMatchedWithinWait_CompetetionStartsAfterNewUserJoin(t *testing.T) {
	setup()
	MatchWaitDuration = 500 * time.Millisecond // Set a short wait duration for testing

	_, err := JoinCompetition("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(1 * time.Second) // Wait for MatchWaitDuration to pass

	comp, err := JoinCompetition("bob_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(1 * time.Second) // Wait for starting competition after MatchWaitDuration

	if comp.StartedAt().IsZero() {
		t.Errorf("competition should have started after %v, got started at %v", MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.Players()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.Players()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
}
