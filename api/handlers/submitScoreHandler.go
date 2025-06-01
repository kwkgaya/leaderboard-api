package handlers

import (
	"encoding/json"
	"net/http"

	"leaderboard/internal/leaderboard"
)

// SubmitScoreHandler godoc
// @Summary      Submit score
// @Description  Add score to the player's current competition
// @Accept       json
// @Param        score  body  map[string]interface{}  true  "Score submission"
// @Success      200  {string}  string  "OK"
// @Failure      409  {string}  string  "Conflict: no active competition"
// @Router       /leaderboard/score [post]
func SubmitScoreHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		PlayerID string `json:"player_id"`
		Score    int    `json:"score"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PlayerID == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}
	if req.Score < 0 {
		http.Error(w, "Score must be a non-negative integer", http.StatusBadRequest)
	}

	// Add score to the player's competing record
	err := leaderboard.AddScore(req.PlayerID, req.Score)
	if err == nil {
		// Successfully added score
		w.WriteHeader(http.StatusOK)
		return
	} else if err == leaderboard.ErrCompetitionEnded {
		http.Error(w, "Competition has ended, cannot add score", http.StatusConflict)
		return
	} else if err == leaderboard.ErrCompetitionNotStarted {
		http.Error(w, "Competition has not started yet, cannot add score", http.StatusConflict)
		return
	} else if err == leaderboard.ErrPlayerNotInCompetition {
		http.Error(w, "Player is not in a competition, cannot add score", http.StatusConflict)
		return
	} else if err == leaderboard.ErrPlayerIdEmpty {
		http.Error(w, "Player ID cannot be empty", http.StatusBadRequest)
		return
	} else if err == leaderboard.ErrPlayerNotFound {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	} else {
		// Some other error occurred
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
