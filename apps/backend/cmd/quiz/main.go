package main

import (
	"log"
	"net/http"

	"github.com/imbivek08/quizz/internal/config"
	"github.com/imbivek08/quizz/internal/handler"
	"github.com/imbivek08/quizz/internal/ws"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize the WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Setup HTTP router
	mux := http.NewServeMux()

	wsHandler := handler.NewWebSocketHandler(hub)
	mux.HandleFunc("/ws", wsHandler.HandleConnection)

	// Health check endpoint
	mux.HandleFunc("/health", handler.HealthCheck)

	// Start server
	log.Printf("🚀 Server starting on %s", cfg.ServerAddress)
	log.Printf("🔌 WebSocket endpoint: ws://localhost%s/ws", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
