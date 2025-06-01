package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// PlayerLeaderboardHandler godoc
// @Summary      Get player leaderboard
// @Description  Get current or last competition for a player
// @Param        playerID  path  string  true  "Player ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /leaderboard/player/{playerID} [get]
func PlayerLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	playerID := chi.URLParam(r, "playerID")

	// Simulate fetching player leaderboard
	response := map[string]string{
		"playerID": playerID,
		"message":  "Player leaderboard fetched successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
