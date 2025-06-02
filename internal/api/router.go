package api

import (
	_ "leaderboard/docs" // Import the generated Swagger docs
	"leaderboard/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Handle("/metrics", promhttp.Handler())

	r.Post("/leaderboard/join", handlers.JoinHandler)
	r.Post("/leaderboard/score", handlers.SubmitScoreHandler)
	r.Get("/leaderboard/player/{playerID}", handlers.PlayerLeaderboardHandler)
	r.Get("/leaderboard/{leaderboardID}", handlers.LeaderboardHandler)

	return r
}
