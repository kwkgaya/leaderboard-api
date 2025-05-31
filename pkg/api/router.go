package api

import (
	"net/http"
	"fmt"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
	_ "leaderboard/docs" // Import the generated Swagger docs
)

func Router() http.Handler {
	r := chi.NewRouter()

	// TODO: Exclude from production builds
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	
	r.Post("/leaderboard/join", joinHandler)
	r.Post("/leaderboard/score", submitScoreHandler)
	r.Get("/leaderboard/player/{playerID}", playerLeaderboardHandler)
	r.Get("/leaderboard/{leaderboardID}", competitionHandler)

	return r
}

// joinHandler godoc
// @Summary      Join a leaderboard competition
// @Description  Match a player to a competition or enqueue them
// @Param        player_id  query  string  true  "Player ID"
// @Success      200  {object}  map[string]interface{}
// @Accepted     202  {string}  string  "Player queued for matchmaking"
// @Failure      400  {string}  string  "Bad request"
// @Failure      409  {string}  string  "Player already in competition"
// @Router       /leaderboard/join [post]
func joinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Simulate joining a leaderboard
	response := map[string]string{
		"message": "Player joined successfully",
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// submitScoreHandler godoc
// @Summary      Submit score
// @Description  Add score to the player's current competition
// @Accept       json
// @Param        score  body  map[string]interface{}  true  "Score submission"
// @Success      200  {string}  string  "OK"
// @Failure      409  {string}  string  "Conflict: no active competition"
// @Router       /leaderboard/score [post]
func submitScoreHandler(w http.ResponseWriter, r *http.Request) {
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

// playerLeaderboardHandler godoc
// @Summary      Get player leaderboard
// @Description  Get current or last competition for a player
// @Param        playerID  path  string  true  "Player ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /leaderboard/player/{playerID} [get]
func playerLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
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

// competitionHandler godoc
// @Summary      Get competition leaderboard
// @Description  Get competition by ID
// @Param        leaderboardID  path  string  true  "Leaderboard ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {string}  string  "Not found"
// @Router       /leaderboard/{leaderboardID} [get]
func competitionHandler(w http.ResponseWriter, r *http.Request) {
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