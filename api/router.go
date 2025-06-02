package api

import (
	"leaderboard/api/handlers"
	_ "leaderboard/docs" // Import the generated Swagger docs
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() http.Handler {
	r := chi.NewRouter()

	// TODO: Exclude from production builds
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/leaderboard/join", handlers.JoinHandler)
	r.Post("/leaderboard/score", handlers.SubmitScoreHandler)
	r.Get("/leaderboard/player/{playerID}", handlers.PlayerLeaderboardHandler)
	r.Get("/leaderboard/{leaderboardID}", handlers.LeaderboardHandler)

	return r
}
