package main


import (
	"log"
	"net/http"
	"leaderboard/pkg/api"
)

func main() {
	server := &http.Server{
		Addr:    ":8080", // TODO: Conmfigure port from environment variable or config file
		Handler: api.Router(),
	}

	log.Println("Starting server on port ", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
