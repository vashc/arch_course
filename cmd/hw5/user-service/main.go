package main

import (
	us "arch_course/internal/hw5/user-service"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create config, e.g. via environment variables
	config, err := us.NewConfig()
	if err != nil {
		log.Fatalf("hw2.NewConfig error: %s", err.Error())
	}

	// Initialize storage
	storage, err := us.NewStorage(config)
	if err != nil {
		log.Fatalf("hw2.NewStorage error: %s", err.Error())
	}

	// Create new fiber application service
	service := us.NewService(config, storage)

	// Instantiate routes
	service.InstantiateRoutes()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down")
	}()

	// Start handling requests
	if err := service.Start(":8000"); err != nil {
		log.Panicf("service.Start error: %s", err.Error())
	}

	log.Println("Running cleanup")
}
