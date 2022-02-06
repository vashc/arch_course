package main

import (
	"arch_course/internal/hw3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create config, e.g. via environment variables
	config, err := hw3.NewConfig()
	if err != nil {
		log.Fatalf("hw3.NewConfig error: %s", err.Error())
	}

	// Initialize storage
	storage, err := hw3.NewStorage(config)
	if err != nil {
		log.Fatalf("hw3.NewStorage error: %s", err.Error())
	}

	// Initialize metrics monitor
	monitor, err := hw3.NewMonitor()
	if err != nil {
		log.Fatalf("hw3.NewMonitor error: %s", err.Error())
	}

	// Create new fiber application service
	service := hw3.NewService(config, storage, monitor)

	// Instantiate routes
	service.InstantiateRoutes()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down")
		_ = service.Shutdown()
	}()

	// Start handling requests
	if err := service.Start(":8000"); err != nil {
		log.Panicf("service.Start error: %s", err.Error())
	}

	log.Println("Running cleanup")
}
