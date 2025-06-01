package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"leaderboard/internal/matchmaking"
	"leaderboard/internal/model"
)

// Save original function to restore after tests
var origJoinCompetition = matchmaking.JoinCompetition

func teardown() {
	matchmaking.JoinCompetition = origJoinCompetition
}

func TestJoinHandler_PlayerIDMissing(t *testing.T) {
	defer teardown()
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Player ID is required") {
		t.Errorf("expected error message for missing player ID, got %s", rr.Body.String())
	}
}

func TestJoinHandler_PlayerIdEmptyError(t *testing.T) {
	defer teardown()
	matchmaking.JoinCompetition = func(playerID string) (model.ICompetition, error) {
		return nil, matchmaking.ErrPlayerIdEmpty
	}
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join?player_id=", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Player ID is required") {
		t.Errorf("expected error message for Player ID is required, got %s", body)
	}
}

func TestJoinHandler_PlayerNotFoundError(t *testing.T) {
	defer teardown()
	matchmaking.JoinCompetition = func(playerID string) (model.ICompetition, error) {
		return nil, matchmaking.ErrPlayerNotFound
	}
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join?player_id=abc", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Player not found") {
		t.Errorf("expected error message for player not found, got %s", rr.Body.String())
	}
}

func TestJoinHandler_PlayerAlreadyInCompetitionError(t *testing.T) {
	defer teardown()
	matchmaking.JoinCompetition = func(playerID string) (model.ICompetition, error) {
		return nil, matchmaking.ErrPlayerAlreadyInCompetition
	}
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join?player_id=abc", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Player already in competition") {
		t.Errorf("expected error message for already in competition	, got %s", rr.Body.String())
	}
}

func TestJoinHandler_InternalServerError(t *testing.T) {
	defer teardown()
	matchmaking.JoinCompetition = func(playerID string) (model.ICompetition, error) {
		return nil, errors.New("unexpected error")
	}
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join?player_id=abc", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Error joining competition") {
		t.Errorf("expected error message for internal error, got %s", rr.Body.String())
	}
}

func TestJoinHandler_PlayerQueued(t *testing.T) {
	defer teardown()
	matchmaking.JoinCompetition = func(playerID string) (model.ICompetition, error) {
		return nil, nil
	}
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join?player_id=abc", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", rr.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["message"] != "Player queued for matchmaking" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
}

func TestJoinHandler_CompetitionStarted(t *testing.T) {
	defer teardown()
	now := time.Now()
	mockComp := &mockCompetition{
		id:        "comp123",
		startedAt: now,
		endsAt:    now.Add(10 * time.Minute),
	}
	matchmaking.JoinCompetition = func(playerID string) (model.ICompetition, error) {
		return mockComp, nil
	}
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join?player_id=abc", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["leaderboard_id"] != "comp123" {
		t.Errorf("unexpected leaderboard_id: %v", resp["leaderboard_id"])
	}
	if int64(resp["ends_at"].(float64)) != mockComp.endsAt.Unix() {
		t.Errorf("unexpected ends_at: %v", resp["ends_at"])
	}
}

func TestJoinHandler_CompetitionNotStarted(t *testing.T) {
	defer teardown()
	mockComp := &mockCompetition{
		id:        "comp123",
		startedAt: time.Time{},
		endsAt:    time.Time{},
	}
	matchmaking.JoinCompetition = func(playerID string) (model.ICompetition, error) {
		return mockComp, nil
	}
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/join?player_id=abc", nil)
	rr := httptest.NewRecorder()

	JoinHandler(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", rr.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["message"] != "Player queued for matchmaking" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
}
