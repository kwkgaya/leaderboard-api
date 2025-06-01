package handlers

import (
	"bytes"
	"errors"
	"io"
	"leaderboard/internal/leaderboard"
	"leaderboard/internal/timeprovider"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

// mock the leaderboard package's GetLeaderboardForPlayer function
var (
	origGetLeaderboardForPlayer = leaderboard.GetLeaderboardForPlayer
)

func mockRouterWithPlayerID(playerID string, _ http.HandlerFunc) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("playerID", playerID)
	req := httptest.NewRequest("GET", "/leaderboard/player/"+playerID, nil)
	return req
}

func TestPlayerLeaderboardHandler_PlayerIdEmpty(t *testing.T) {
	leaderboard.GetLeaderboardForPlayer = func(playerID string) (*leaderboard.LeaderboardResponse, error) {
		return nil, leaderboard.ErrPlayerIdEmpty
	}
	defer func() { leaderboard.GetLeaderboardForPlayer = origGetLeaderboardForPlayer }()

	req := mockRouterWithPlayerID("", PlayerLeaderboardHandler)
	rr := httptest.NewRecorder()

	PlayerLeaderboardHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if !bytes.Contains(body, []byte("Player ID cannot be empty")) {
		t.Errorf("expected error message, got %s", string(body))
	}
}

func TestPlayerLeaderboardHandler_PlayerNotFound(t *testing.T) {
	leaderboard.GetLeaderboardForPlayer = func(playerID string) (*leaderboard.LeaderboardResponse, error) {
		return nil, leaderboard.ErrPlayerNotFound
	}
	defer func() { leaderboard.GetLeaderboardForPlayer = origGetLeaderboardForPlayer }()

	req := mockRouterWithPlayerID("notfound", PlayerLeaderboardHandler)
	rr := httptest.NewRecorder()

	PlayerLeaderboardHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if !bytes.Contains(body, []byte("Player not found")) {
		t.Errorf("expected error message, got %s", string(body))
	}
}

func TestPlayerLeaderboardHandler_Success(t *testing.T) {
	expected := &leaderboard.LeaderboardResponse{
		Id:     "leaderboard123",
		EndsAt: timeprovider.Current.Now().Add(1 * time.Hour),
		Leaderboard: []leaderboard.PlayerScore{
			{PlayerId: "123", Score: 100},
			{PlayerId: "456", Score: 200}}}

	leaderboard.GetLeaderboardForPlayer = func(playerID string) (*leaderboard.LeaderboardResponse, error) {
		return expected, nil
	}
	defer func() { leaderboard.GetLeaderboardForPlayer = origGetLeaderboardForPlayer }()

	req := mockRouterWithPlayerID("123", PlayerLeaderboardHandler)
	rr := httptest.NewRecorder()

	PlayerLeaderboardHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", rr.Header().Get("Content-Type"))
	}
	body, _ := io.ReadAll(rr.Body)
	if !bytes.Contains(body, []byte(`"score":100`)) {
		t.Errorf("expected response body to contain score, got %s", string(body))
	}
}

func TestPlayerLeaderboardHandler_UnexpectedError(t *testing.T) {
	leaderboard.GetLeaderboardForPlayer = func(playerID string) (*leaderboard.LeaderboardResponse, error) {
		return nil, errors.New("unexpected error")
	}
	defer func() { leaderboard.GetLeaderboardForPlayer = origGetLeaderboardForPlayer }()

	req := mockRouterWithPlayerID("123", PlayerLeaderboardHandler)
	rr := httptest.NewRecorder()

	PlayerLeaderboardHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
}
