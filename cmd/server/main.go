package main

import (
	"fmt"
	"log"

	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Starting server at %s:%s\n", cfg.Server.Host, cfg.Server.Port)

	// Create server (now returns error following Go philosophy)
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
