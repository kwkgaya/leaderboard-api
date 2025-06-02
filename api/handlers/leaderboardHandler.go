package handlers

import (
	"encoding/json"
	"fmt"
	"leaderboard/internal/leaderboard"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// LeaderboardHandler godoc
// @Summary      Get leaderboard
// @Description  Get leaderboard by ID
// @Param        leaderboardID  path  string  true  "Leaderboard ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {string}  string  "Not found"
// @Router       /leaderboard/{leaderboardID} [get]
func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {

	leaderboardID := chi.URLParam(r, "leaderboardID")

	response, err := leaderboard.GetLeaderboard(leaderboardID)
	if err == leaderboard.ErrCompetetionNotFound {
		http.Error(w, "Leaderboard not found", http.StatusNotFound)
		return
	} else if err == leaderboard.ErrLeaderboardIdEmpty {
		http.Error(w, "Leaderboard ID cannot be empty", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	} else if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}
}
