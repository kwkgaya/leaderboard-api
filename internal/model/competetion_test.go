package model

import (
	"fmt"
	"testing"
	"time"

	"leaderboard/internal/config"
	"leaderboard/internal/timeprovider"

	"github.com/google/uuid"
)

func TestNewCompetition_InitializesFieldsCorrectly(t *testing.T) {
	// Arrange
	fixedTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	originalProvider := timeprovider.Current
	timeprovider.Current = &timeprovider.MockTimeProvider{FixedTime: fixedTime}
	defer func() { timeprovider.Current = originalProvider }()

	// Act
	competition := NewCompetition()

	// Assert
	if competition == nil {
		t.Fatal("expected non-nil Competition")
	}
	if _, err := uuid.Parse(competition.Id()); err != nil {
		t.Errorf("expected valid UUID for ID, got %q", competition.Id())
	}
	if !competition.CreatedAt().Equal(fixedTime) {
		t.Errorf("expected CreatedAt %v, got %v", fixedTime, competition.CreatedAt())
	}
	if !competition.StartedAt().IsZero() {
		t.Errorf("expected StartedAt to be zero, got %v", competition.StartedAt())
	}
	if !competition.EndsAt().IsZero() {
		t.Errorf("expected EndsAt to be zero, got %v", competition.EndsAt())
	}
	if competition.Players() == nil {
		t.Error("expected Players slice to be initialized, got nil")
	}
	if cap(competition.Players()) != config.MaxPlayersForCompetetion {
		t.Errorf("expected Players cap %d, got %d", config.MaxPlayersForCompetetion, cap(competition.Players()))
	}
	if len(competition.Players()) != 0 {
		t.Errorf("expected Players len 0, got %d", len(competition.Players()))
	}
}

func TestCompetition_AddPlayer_Success(t *testing.T) {
	competition := NewCompetition()
	player := NewPlayer("p1", 1, "US")

	err := competition.AddPlayer(player)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(competition.Players()) != 1 {
		t.Errorf("expected 1 player, got %d", len(competition.Players()))
	}
	if competition.Players()[0].Player() != player {
		t.Error("expected player to be added to Players slice")
	}
	if player.ActiveCompetition() != competition {
		t.Error("expected player's ActiveCompetition to be set")
	}
	// Score should be 0
	if competition.Players()[0].Score() != 0 {
		t.Errorf("expected score 0, got %d", competition.Players()[0].Score())
	}
}

func TestCompetition_AddPlayer_CompetitionFull(t *testing.T) {
	competition := NewCompetition()
	for i := 0; i < config.MaxPlayersForCompetetion; i++ {
		player := NewPlayer(fmt.Sprintf("p%v", i+1), 1, "US")

		err := competition.AddPlayer(player)
		if err != nil {
			t.Fatalf("unexpected error adding player %d: %v", i, err)
		}
	}
	player := NewPlayer("overflow", 1, "US")
	err := competition.AddPlayer(player)
	if err != ErrCompetitionFull {
		t.Errorf("expected ErrCompetitionFull, got %v", err)
	}
}

func TestCompetition_AddPlayer_CompetitionStarted(t *testing.T) {
	competition := NewCompetition()
	for i := 0; i < config.MinPlayersForCompetetion; i++ {
		err1 := competition.AddPlayer(NewPlayer(fmt.Sprintf("p%v", i+1), 1, "US"))
		if err1 != nil {
			t.Fatalf("competetion.AddPlayer() returned error %v", err1)
		}
	}
	err2 := competition.Start()
	if err2 != nil {
		t.Fatalf("competetion.Start() returned error %v", err2)
	}

	player := NewPlayer("p0", 1, "US")
	err := competition.AddPlayer(player)
	if err != ErrCompetitionStarted {
		t.Errorf("expected ErrCompetitionStarted, got %v", err)
	}
}

func TestCompetition_Start_NotEnoughPlayers(t *testing.T) {
	competition := NewCompetition()
	// Only one player, less than MinPlayersForCompetetion
	player := NewPlayer("p0", 1, "US")
	_ = competition.AddPlayer(player)

	err := competition.Start()

	if err != ErrNotEnoughPlayers {
		t.Errorf("expected ErrNotEnoughPlayers, got %v", err)
	}
}

func TestCompetition_Start_SuccessWithMinPlayers(t *testing.T) {
	fixedTime := time.Date(2024, 6, 3, 15, 0, 0, 0, time.UTC)
	originalProvider := timeprovider.Current
	timeprovider.Current = &timeprovider.MockTimeProvider{FixedTime: fixedTime}
	defer func() { timeprovider.Current = originalProvider }()

	competition := NewCompetition()
	player1 := NewPlayer("p1", 1, "US")
	player2 := NewPlayer("p2", 1, "US")
	_ = competition.AddPlayer(player1)
	_ = competition.AddPlayer(player2)

	err := competition.Start()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !competition.StartedAt().Equal(fixedTime) {
		t.Errorf("expected StartedAt %v, got %v", fixedTime, competition.StartedAt())
	}
	expectedEndsAt := fixedTime.Add(config.CompetitionDuration)
	if !competition.EndsAt().Equal(expectedEndsAt) {
		t.Errorf("expected EndsAt %v, got %v", expectedEndsAt, competition.EndsAt())
	}
}

func TestCompetition_Start_CalledTwice(t *testing.T) {
	competition := NewCompetition()
	player1 := NewPlayer("p1", 1, "US")
	player2 := NewPlayer("p2", 1, "US")
	_ = competition.AddPlayer(player1)
	_ = competition.AddPlayer(player2)

	err1 := competition.Start()
	err2 := competition.Start()

	if err1 != nil {
		t.Errorf("expected first Start to succeed, got %v", err1)
	}
	if err2 != ErrCompetitionStarted {
		t.Errorf("expected ErrCompetitionStarted on second Start, got %v", err2)
	}
}
