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

// mockGetLeaderboard is used to mock leaderboard.GetLeaderboard in tests
var mockGetLeaderboard func(string) (*leaderboard.LeaderboardResponse, error)

func setupMock() func() {
	orig := leaderboard.GetLeaderboard
	leaderboard.GetLeaderboard = func(id string) (*leaderboard.LeaderboardResponse, error) {
		return mockGetLeaderboard(id)
	}
	return func() { leaderboard.GetLeaderboard = orig }
}

func TestLeaderboardHandler_Success(t *testing.T) {
	restore := setupMock()
	defer restore()

	expected := &leaderboard.LeaderboardResponse{
		Id:     "leaderboard123",
		EndsAt: timeprovider.Current.Now().Add(1 * time.Hour),
		Leaderboard: []leaderboard.PlayerScore{
			{PlayerId: "123", Score: 100},
			{PlayerId: "456", Score: 200}}}

	mockGetLeaderboard = func(id string) (*leaderboard.LeaderboardResponse, error) {
		return expected, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/leaderboard/leaderboard123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("leaderboardID", "leaderboard123")

	rr := httptest.NewRecorder()
	LeaderboardHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte(`"leaderboard_id":"leaderboard123"`)) {
		t.Errorf("expected response body to contain leaderboard id, got %s", string(body))
	}
}

func TestLeaderboardHandler_NotFound(t *testing.T) {
	restore := setupMock()
	defer restore()

	mockGetLeaderboard = func(id string) (*leaderboard.LeaderboardResponse, error) {
		return nil, leaderboard.ErrCompetetionNotFound
	}

	req := httptest.NewRequest(http.MethodGet, "/leaderboard/notfound", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("leaderboardID", "notfound")

	rr := httptest.NewRecorder()
	LeaderboardHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("Leaderboard not found")) {
		t.Errorf("expected not found message, got %s", string(body))
	}
}

func TestLeaderboardHandler_EmptyID(t *testing.T) {
	restore := setupMock()
	defer restore()

	mockGetLeaderboard = func(id string) (*leaderboard.LeaderboardResponse, error) {
		return nil, leaderboard.ErrLeaderboardIdEmpty
	}

	req := httptest.NewRequest(http.MethodGet, "/leaderboard/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("leaderboardID", "")

	rr := httptest.NewRecorder()
	LeaderboardHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("Leaderboard ID cannot be empty")) {
		t.Errorf("expected empty id message, got %s", string(body))
	}
}

func TestLeaderboardHandler_InternalError(t *testing.T) {
	restore := setupMock()
	defer restore()

	mockGetLeaderboard = func(_ string) (*leaderboard.LeaderboardResponse, error) {
		return nil, errors.New("db error")
	}

	req := httptest.NewRequest(http.MethodGet, "/leaderboard/abc", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("leaderboardID", "abc")

	rr := httptest.NewRecorder()
	LeaderboardHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("Internal server error")) {
		t.Errorf("expected internal error message, got %s", string(body))
	}
}
