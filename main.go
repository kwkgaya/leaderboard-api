package main


import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"

	"leaderboard/pkg/api"
)

func main() {

	server := &http.Server{
		Addr:    ":8080", // TODO: Conmfigure port from environment variable or config file
		Handler: api.Router(),
	}

	go func() {
		log.Println("Starting server on port ", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
		
		log.Println("Server stopped gracefully")
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: Configure graceful shutdown timeout
	defer cancel()
	_ = server.Shutdown(ctx)
}
