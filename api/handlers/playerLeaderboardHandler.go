package handlers

import (
	"encoding/json"
	"fmt"
	"leaderboard/internal/leaderboard"
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

	playerID := chi.URLParam(r, "playerID")

	response, err := leaderboard.GetLeaderboardForPlayer(playerID)
	if err == leaderboard.ErrPlayerIdEmpty {
		http.Error(w, "Player ID cannot be empty", http.StatusBadRequest)
		return
	} else if err == leaderboard.ErrPlayerNotFound {
		http.Error(w, "Player not found", http.StatusBadRequest)
		return
	} else if err == leaderboard.ErrPlayerNotInCompetition {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	} else if encodingErr := json.NewEncoder(w).Encode(response); encodingErr != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", encodingErr), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}
}
