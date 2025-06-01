package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CompetitionHandler godoc
// @Summary      Get competition leaderboard
// @Description  Get competition by ID
// @Param        leaderboardID  path  string  true  "Leaderboard ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {string}  string  "Not found"
// @Router       /leaderboard/{leaderboardID} [get]
func CompetitionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	leaderboardID := chi.URLParam(r, "leaderboardID")

	// Simulate fetching competition leaderboard
	response := map[string]string{
		"leaderboardID": leaderboardID,
		"message":       "Competition leaderboard fetched successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
