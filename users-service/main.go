package main

import (
	"log"

	"github.com/oyinetare/go-docker-microservice/config"
	"github.com/oyinetare/go-docker-microservice/repository"
	"github.com/oyinetare/go-docker-microservice/server"
)

func main() {
	log.Println("--- Customer Service ---")
	log.Println("Connecting to customer repository...")

	// load config
	cfg := config.New()

	// connect to db
	repo, err := repository.Connect(
		cfg.DB.Host,
		cfg.DB.Database,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Port,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Disconnect()

	log.Println("Connected. Starting server...")

	// create and start server
	srv := server.New(repo, cfg.Port)
	log.Printf("Server started successfully, running on port %d", cfg.Port)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
