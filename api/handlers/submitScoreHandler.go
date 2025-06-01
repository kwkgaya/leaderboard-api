package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")

	// Simulate score submission
	response := map[string]string{
		"message": "Score submitted successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
