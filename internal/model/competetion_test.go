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
	if competition.PlayersMap() == nil {
		t.Error("expected Players map to be initialized, got nil")
	}
	if len(competition.PlayersMap()) != 0 {
		t.Errorf("expected Players map to be empty, got %d", len(competition.PlayersMap()))
	}
	if len(competition.PlayersMap()) != 0 {
		t.Errorf("expected Players map len 0, got %d", len(competition.PlayersMap()))
	}
}

func TestCompetition_AddPlayer_Success(t *testing.T) {
	competition := NewCompetition()
	player := NewPlayer("p1", 1, "US")

	err := competition.AddPlayer(player)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(competition.PlayersMap()) != 1 {
		t.Errorf("expected 1 player, got %d", len(competition.PlayersMap()))
	}
	compPlayer, ok := competition.PlayersMap()[player.Id()]
	if !ok || compPlayer.Player() != player {
		t.Error("expected player to be added to Players map")
	}
	if player.Competition() != competition {
		t.Error("expected player's Competition to be set")
	}
	// Score should be 0
	compPlayer, ok = competition.PlayersMap()[player.Id()]
	if !ok || compPlayer.Score() != 0 {
		t.Errorf("expected score 0, got %d", compPlayer.Score())
	}
}

func TestCompetition_AddPlayer_CompetitionFull(t *testing.T) {
	competition := NewCompetition()
	for i := 0; i < config.MaxPlayersForCompetition; i++ {
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
	for i := 0; i < config.MinPlayersForCompetition; i++ {
		err1 := competition.AddPlayer(NewPlayer(fmt.Sprintf("p%v", i+1), 1, "US"))
		if err1 != nil {
			t.Fatalf("competition.AddPlayer() returned error %v", err1)
		}
	}
	err2 := competition.Start()
	if err2 != nil {
		t.Fatalf("competition.Start() returned error %v", err2)
	}

	player := NewPlayer("p0", 1, "US")
	err := competition.AddPlayer(player)
	if err != ErrCompetitionStarted {
		t.Errorf("expected ErrCompetitionStarted, got %v", err)
	}
}

func TestCompetition_Start_NotEnoughPlayers(t *testing.T) {
	competition := NewCompetition()
	// Only one player, less than MinPlayersForCompetition
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

func TestCompetition_AddScore_Success(t *testing.T) {
	competition := NewCompetition()
	player := NewPlayer("p1", 1, "US")
	err := competition.AddPlayer(player)
	if err != nil {
		t.Fatalf("unexpected error adding player: %v", err)
	}

	err = competition.AddScore(player, 10)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	compPlayer, ok := competition.PlayersMap()[player.Id()]
	if !ok {
		t.Fatalf("player not found in competition after AddScore")
	}
	if compPlayer.Score() != 10 {
		t.Errorf("expected score 10, got %d", compPlayer.Score())
	}
}

func TestCompetition_AddScore_PlayerNil(t *testing.T) {
	competition := NewCompetition()
	err := competition.AddScore(nil, 5)
	if err != ErrPlayerNil {
		t.Errorf("expected ErrPlayerNil, got %v", err)
	}
}

func TestCompetition_AddScore_PointsNegative(t *testing.T) {
	competition := NewCompetition()
	player := NewPlayer("p1", 1, "US")
	_ = competition.AddPlayer(player)

	err := competition.AddScore(player, -1)
	if err != ErrPointsNegative {
		t.Errorf("expected ErrPointsNegative, got %v", err)
	}
}

func TestCompetition_AddScore_PlayerNotFound(t *testing.T) {
	competition := NewCompetition()
	player := NewPlayer("p1", 1, "US")
	// Not adding player to competition

	err := competition.AddScore(player, 5)
	if err != ErrPlayerNotFound {
		t.Errorf("expected ErrPlayerNotFound, got %v", err)
	}
}

func TestCompetition_AddScore_SortLeaderboardAccordingToScore(t *testing.T) {
	competition := NewCompetition()
	player1 := NewPlayer("a", 1, "US")
	player2 := NewPlayer("b", 1, "US")
	_ = competition.AddPlayer(player1)
	_ = competition.AddPlayer(player2)
	_ = competition.Start()

	_ = competition.AddScore(player1, 10)
	_ = competition.AddScore(player2, 20)

	leaderboard := competition.Leaderboard()
	if len(leaderboard) != 2 {
		t.Fatalf("expected 2 sorted players, got %d", len(leaderboard))
	}
	if leaderboard[0].Player().Id() != "b" || leaderboard[1].Player().Id() != "a" {
		t.Errorf("expected Leaderboard()[0] = b, Leaderboard()[1] = a")
	}
	// Now add more score to player1 so they overtake player2
	_ = competition.AddScore(player1, 15) // player1 now has 25

	leaderboard = competition.Leaderboard()
	if leaderboard[0].Player().Id() != "a" || leaderboard[1].Player().Id() != "b" {
		t.Errorf("expected Leaderboard()[0] = a, Leaderboard()[1] = b after score update")
	}
}

func TestCompetition_AddScore_SortLeaderboardAccordingToScoreThenName(t *testing.T) {
	competition := NewCompetition()
	player1 := NewPlayer("a", 1, "US")
	player2 := NewPlayer("b", 1, "US")
	player3 := NewPlayer("c", 1, "US")
	_ = competition.AddPlayer(player1)
	_ = competition.AddPlayer(player2)
	_ = competition.AddPlayer(player3)
	_ = competition.Start()

	_ = competition.AddScore(player1, 10)
	_ = competition.AddScore(player2, 20)
	_ = competition.AddScore(player3, 30)

	leaderboard := competition.Leaderboard()
	if len(leaderboard) != 3 {
		t.Fatalf("expected 2 sorted players, got %d", len(leaderboard))
	}
	if leaderboard[0].Player().Id() != "c" || leaderboard[1].Player().Id() != "b" || leaderboard[2].Player().Id() != "a" {
		t.Errorf("expected Leaderboard()[0] = c, Leaderboard()[1] = b, Leaderboard()[2] = a")
	}
	// Now add more score to player1 so they are equal to player3
	_ = competition.AddScore(player1, 20)
	leaderboard = competition.Leaderboard()

	if leaderboard[0].Player().Id() != "a" || leaderboard[1].Player().Id() != "c" || leaderboard[2].Player().Id() != "b" {
		t.Errorf("expected Leaderboard()[0] = a, Leaderboard()[1] = c, Leaderboard()[2] = b after score update")
	}

	_ = competition.AddScore(player2, 20)
	_ = competition.AddScore(player3, 10)
	leaderboard = competition.Leaderboard()

	if leaderboard[0].Player().Id() != "b" || leaderboard[1].Player().Id() != "c" || leaderboard[2].Player().Id() != "a" {
		t.Errorf("expected Leaderboard()[0] = b, Leaderboard()[1] = c, Leaderboard()[2] = a after score update")
	}
}
