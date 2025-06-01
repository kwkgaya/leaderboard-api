package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"leaderboard/internal/leaderboard"
)

// Mock leaderboard.AddScore and error variables for testing
var (
	mockAddScoreFunc func(playerID string, score int) error
)

func mockAddScore(playerID string, score int) error {
	return mockAddScoreFunc(playerID, score)
}

func setupMocks() func() {
	origAddScore := leaderboard.AddScore
	leaderboard.AddScore = mockAddScore
	return func() { leaderboard.AddScore = origAddScore }
}

func TestSubmitScoreHandler_Success(t *testing.T) {
	restore := setupMocks()
	defer restore()
	mockAddScoreFunc = func(playerID string, score int) error {
		return nil
	}

	body := []byte(`{"player_id":"player1","score":100}`)
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/score", bytes.NewReader(body))
	w := httptest.NewRecorder()

	SubmitScoreHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestSubmitScoreHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/score", bytes.NewReader([]byte("{invalid json")))
	w := httptest.NewRecorder()

	SubmitScoreHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Invalid request body\n" {
		t.Errorf("unexpected body: %s", string(body))
	}
}

func TestSubmitScoreHandler_EmptyPlayerID(t *testing.T) {
	body := []byte(`{"player_id":"","score":10}`)
	req := httptest.NewRequest(http.MethodPost, "/leaderboard/score", bytes.NewReader(body))
	w := httptest.NewRecorder()

	SubmitScoreHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestSubmitScoreHandler_ErrorCases_TableDriven(t *testing.T) {
	restore := setupMocks()
	defer restore()

	tests := []struct {
		name           string
		errorToReturn  error
		expectedStatus int
	}{
		{
			name:           "CompetitionEnded",
			errorToReturn:  leaderboard.ErrCompetitionEnded,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "CompetitionNotStarted",
			errorToReturn:  leaderboard.ErrCompetitionNotStarted,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "PlayerNotInCompetition",
			errorToReturn:  leaderboard.ErrPlayerNotInCompetition,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "PlayerIdEmptyError",
			errorToReturn:  leaderboard.ErrPlayerIdEmpty,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "PlayerNotFound",
			errorToReturn:  leaderboard.ErrPlayerNotFound,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InternalServerError",
			errorToReturn:  errors.New("some internal error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	body := []byte(`{"player_id":"player1","score":10}`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAddScoreFunc = func(_ string, _ int) error {
				return tt.errorToReturn
			}
			req := httptest.NewRequest(http.MethodPost, "/leaderboard/score", bytes.NewReader(body))
			w := httptest.NewRecorder()

			SubmitScoreHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
