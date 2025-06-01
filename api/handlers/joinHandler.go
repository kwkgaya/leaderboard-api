package handlers

import (
	"encoding/json"
	"fmt"
	"leaderboard/internal/matchmaking"
	"net/http"
)

// JoinHandler godoc
// @Summary      Join a leaderboard competition
// @Description  Match a player to a competition or enqueue them
// @Param        player_id  query  string  true  "Player ID"
// @Success      200  {object}  map[string]interface{}
// @Accepted     202  {string}  string  "Player queued for matchmaking"
// @Failure      400  {string}  string  "Player ID is empty or player not found"
// @Failure      409  {string}  string  "Player already in competition"
// @Router       /leaderboard/join [post]
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	comp, err := matchmaking.JoinCompetition(playerID)
	if err != nil {
		if err == matchmaking.ErrPlayerIdEmpty {
			http.Error(w, "Player ID cannot be empty", http.StatusBadRequest)
			return
		} else if err == matchmaking.ErrPlayerNotFound {
			http.Error(w, "Player not found", http.StatusBadRequest)
			return
		} else if err == matchmaking.ErrPlayerAlreadyInCompetition {
			http.Error(w, "Player already in competition", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Error joining competition: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// If comp is nil, player is queued for matchmaking
	if comp == nil {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Player queued for matchmaking",
		})
		return
	}

	// If competition has started, return 200 with leaderboard_id and ends_at
	if !comp.StartedAt().IsZero() {
		resp := map[string]interface{}{
			"leaderboard_id": comp.Id(),
			"ends_at":        comp.EndsAt().Unix(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// If competition exists but hasn't started, player is still queued
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Player queued for matchmaking",
	})
}
