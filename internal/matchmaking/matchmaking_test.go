package matchmaking

import (
	"errors"
	"fmt"
	"leaderboard/internal/config"
	"leaderboard/internal/model"
	"leaderboard/internal/storage"
	"testing"
	"time"
)

func setup() {
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

func tearDown() {
	config.MatchWaitDuration = 30 * time.Second
	config.CompetitionDuration = 1 * time.Hour

	clear(waitingCompetitions)
	clear(storage.Players)
	clear(storage.Competitions)
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
	tearDown()
}

func TestJoinCompetition_Matchmaking(t *testing.T) {
	setup()

	// First player joins, should be put in waiting list
	comp1, err := JoinCompetition("bob")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp1 == nil {
		t.Fatalf("expected competition to be created for bob")
	}
	if waitingCompetitions[2] == nil || waitingCompetitions[2].PlayersMap()["bob"] == nil {
		t.Errorf("waitingCompetitions at level 2 should have bob, got %v", waitingCompetitions[2])
	}

	// Second player joins, should create a competition
	comp2, err := JoinCompetition("bob_1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp2 == nil {
		t.Fatalf("expected competition to be created for bob_1")
	}
	if comp1.Id() != comp2.Id() {
		t.Errorf("expected same competition for bob and bob_1, got different competitions %s and %s", comp1.Id(), comp2.Id())
	}
	if len(comp1.PlayersMap()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp1.PlayersMap()))
	}
	if !comp1.StartedAt().IsZero() {
		t.Errorf("competition should not have started yet, got started at %v", comp1.StartedAt())
	}
	if waitingCompetitions[2] == nil {
		t.Errorf("waitingCompetitions at level 2 should not be nil, got %v", waitingCompetitions)
	}
	tearDown()
}

func TestJoinCompetition_AlreadyInCompetition(t *testing.T) {
	setup()

	storage.AddPlayers([]storage.NewPlayer{
		{Id: "player3", CountryCode: "US", Level: 7},
	})

	player := storage.Players["player3"]
	// Simulate player already in a competition
	fakeComp := model.NewCompetition(1)
	player.SetCompetition(fakeComp)

	comp, err := JoinCompetition("player3")

	if !errors.Is(err, ErrPlayerAlreadyInCompetition) {
		t.Errorf("expected ErrPlayerAlreadyInCompetition, got %v", err)
	}
	if comp != nil {
		t.Errorf("expected nil competition, got %v", comp)
	}
	tearDown()
}

func TestJoinCompetition_JoinMaxplayers_CompetitionStarts(t *testing.T) {
	setup()
	var previousComp model.ICompetition
	for i := 1; i <= config.MaxPlayersForCompetition; i++ {
		playerId := fmt.Sprintf("player%v", i)
		storage.AddPlayers([]storage.NewPlayer{
			{Id: playerId, CountryCode: "US", Level: 5},
		})

		comp, err := JoinCompetition(playerId)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if comp == nil {
			t.Fatalf("expected competition to be returned for %s", playerId)
		}
		if len(comp.PlayersMap()) != i {
			t.Errorf("competition should have %d players, got %d", i, len(comp.PlayersMap()))
		}
		if previousComp != nil && comp.Id() != previousComp.Id() {
			t.Errorf("expected same competition for player %s, got different competition %s", playerId, comp.Id())
		}
		previousComp = comp

		player := storage.Players[playerId]
		if player.Competition() == nil || player.Competition().Id() != comp.Id() {
			t.Errorf("player %s should be in competition %s, got %v", playerId, comp.Id(), player.Competition())
		}

		if i == config.MaxPlayersForCompetition {
			if comp.StartedAt().IsZero() {
				t.Errorf("competition should have started after adding %d players, got started at %v", config.MaxPlayersForCompetition, comp.StartedAt())
			}
			if len(waitingCompetitions) != 0 {
				t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
			}

			// Test player1 only one time
			player1Id := fmt.Sprintf("player%v", 1)
			player1 := storage.Players[player1Id]
			if player1.Competition() == nil || player1.Competition().Id() != comp.Id() {
				t.Errorf("player %s should be in competition %s, got %v", player1Id, comp.Id(), player1.Competition())
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
	tearDown()
}

func TestJoinCompetition_JoinMaxPlayersAndTwoMore_NewCompetitionCreated(t *testing.T) {
	setup()
	var previousComp model.ICompetition
	for i := 1; i <= config.MaxPlayersForCompetition+2; i++ {
		playerId := fmt.Sprintf("player%v", i)
		storage.AddPlayers([]storage.NewPlayer{
			{Id: playerId, CountryCode: "US", Level: 5},
		})

		comp, err := JoinCompetition(playerId)

		if i == config.MaxPlayersForCompetition {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			previousComp = comp
		} else if i > config.MaxPlayersForCompetition+2 {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if comp.Id() == previousComp.Id() {
				t.Fatalf("expected new competition to be created for %s, got same competition %s", playerId, comp.Id())
			}
		}
	}
	tearDown()
}

func TestJoinCompetition_MatchedwithTwoPlayers_CompetitionStartsAfterMatchWaitDuration(t *testing.T) {
	setup()
	config.MatchWaitDuration = 500 * time.Millisecond // Set a short wait duration for testing

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
		t.Errorf("competition should have started after %v, got started at %v", config.MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.PlayersMap()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.PlayersMap()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
	tearDown()
}

func TestJoinCompetition_MatchedWithTwoPlayersInTwoLevels_CompetitionStartsAfterMatchWaitDuration(t *testing.T) {
	setup()
	config.MatchWaitDuration = 1 * time.Second // Set a short wait duration for testing

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

	if alice.Competition() == nil {
		t.Errorf("alice should be in a competition, got nil")
	}
	if bob.Competition() == nil {
		t.Errorf("bob should be in a competition, got nil")
	}
	if alice.Competition().Id() != bob.Competition().Id() {
		t.Errorf("alice and bob should be in the same competition, got %s and %s", alice.Competition().Id(), bob.Competition().Id())
	}
	comp := alice.Competition()

	if comp.StartedAt().IsZero() {
		t.Errorf("competition should have started after %v, got started at %v", config.MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.PlayersMap()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.PlayersMap()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
	tearDown()
}

func TestJoinCompetition_MatchedwithTwoPlayersInMinAndMaxLevels_CompetitionStartsAfterMatchWaitDuration(t *testing.T) {
	setup()
	config.MatchWaitDuration = 1 * time.Second // Set a short wait duration for testing

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

	if alice.Competition() == nil {
		t.Errorf("alice should be in a competition, got nil")
	}
	if ian.Competition() == nil {
		t.Errorf("ian should be in a competition, got nil")
	}
	if alice.Competition().Id() != ian.Competition().Id() {
		t.Errorf("alice and bob should be in the same competition, got %s and %s", alice.Competition().Id(), ian.Competition().Id())
	}
	comp := alice.Competition()

	if comp.StartedAt().IsZero() {
		t.Errorf("competition should have started after %v, got started at %v", config.MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.PlayersMap()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.PlayersMap()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
	tearDown()
}

func TestJoinCompetition_CompetitionStartAfterWait_NewJoineesAddedToNewCompetition(t *testing.T) {
	setup()
	config.MatchWaitDuration = 1500 * time.Millisecond // Set a short wait duration for testing

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
	comp1 := alice.Competition()

	alice1comp, err := JoinCompetition("alice_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alice1comp == nil {
		t.Fatalf("expected alice_1 to be added to competition, got nil")
	}

	bob1comp, err := JoinCompetition("bob_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bob1comp == nil {
		t.Fatalf("expected bob_1 to be added to competition, got nil")
	}
	if waitingCompetitions[2] == nil {
		t.Errorf("waitingCompetitions at level 2 should not be nil, got %v", waitingCompetitions)
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
	tearDown()
}

func TestJoinCompetition_NotMatchedWithinWait_CompetitionStartsAfterNewUserJoin(t *testing.T) {
	setup()
	config.MatchWaitDuration = 500 * time.Millisecond // Set a short wait duration for testing

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
		t.Errorf("competition should have started after %v, got started at %v", config.MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.PlayersMap()) != 2 {
		t.Errorf("competition should have 2 players, got %d", len(comp.PlayersMap()))
	}
	if len(comp.Leaderboard()) != 2 {
		t.Errorf("competition should have 2 players in leaderboard, got %d", len(comp.Leaderboard()))
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
	tearDown()
}

func TestJoinCompetition_CompetetionStartAndEnd_UserCanJoinNewCompetition(t *testing.T) {
	setup()
	config.MatchWaitDuration = 500 * time.Millisecond // Set a short wait duration for testing
	config.CompetitionDuration = 1 * time.Second

	_, err := JoinCompetition("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	comp1, err := JoinCompetition("bob_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(2 * time.Second) // Wait for starting and ending

	_, _ = JoinCompetition("bob_2")
	comp2, err := JoinCompetition("bob_1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp1 == nil {
		t.Fatalf("expected competition 1 to be created for bob_1, got nil")
	}
	if comp2 == nil {
		t.Fatalf("expected new competition 2 to be created for bob_1, got nil")
	}
	if comp1.Id() == comp2.Id() {
		t.Errorf("expected different competitions for bob_1, got same competition %s", comp1.Id())
	}
	tearDown()
}

func TestJoinCompetition_CompetitionCreatedAtAdjacentLevel_CompetitionStartsAfterMatchWaitDuration(t *testing.T) {
	setup()
	config.MatchWaitDuration = 1 * time.Second // Set a short wait duration for testing

	aliceComp, err := JoinCompetition("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bobComp, err := JoinCompetition("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bob1Comp, err := JoinCompetition("bob_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if aliceComp == bobComp {
		t.Errorf("expected alice and bob to be in different competitions, got same competition %s", aliceComp.Id())
	}
	if bobComp != bob1Comp {
		t.Errorf("expected bob and bob_1 to be in the same competition, got different competitions %s and %s", bobComp.Id(), bob1Comp.Id())
	}

	alice := storage.Players["alice"]
	bob := storage.Players["bob"]
	bob1 := storage.Players["bob_1"]

	time.Sleep(2 * time.Second) // Wait for starting competition after MatchWaitDuration

	if alice.Competition() == nil {
		t.Errorf("alice should be in a competition, got nil")
	}
	if bob.Competition() == nil {
		t.Errorf("bob should be in a competition, got nil")
	}
	if bob1.Competition() == nil {
		t.Errorf("bob_1 should be in a competition, got nil")
	}
	if bob.Competition().Id() != bob1.Competition().Id() {
		t.Errorf("bob and bob_1 should be in the same competition, got %s and %s", bob.Competition().Id(), bob1.Competition().Id())
	}
	if alice.Competition().Id() != bob.Competition().Id() {
		t.Errorf("alice and bob should be in the same competition, got %s and %s", alice.Competition().Id(), bob.Competition().Id())
	}
	comp := alice.Competition()
	if comp != bobComp {
		t.Errorf("expected alice's competition to be the same as bob's result from JoinCompetetion, got %s and %s", comp.Id(), bobComp.Id())
	}

	if comp.StartedAt().IsZero() {
		t.Errorf("competition should have started after %v, got started at %v", config.MatchWaitDuration, comp.StartedAt())
	}
	if len(comp.PlayersMap()) != 3 {
		t.Errorf("competition should have 3 players, got %v", comp.PlayersMap())
	}
	if len(waitingCompetitions) != 0 {
		t.Errorf("waitingCompetitions should be empty after competition started, got %v", waitingCompetitions)
	}
	tearDown()
}
