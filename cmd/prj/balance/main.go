package main

import (
	balance "arch_course/internal/prj/balance"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create config, e.g. via environment variables
	config, err := balance.NewConfig()
	if err != nil {
		log.Fatalf("balance.NewConfig error: %s", err.Error())
	}

	// Initialize storage
	storage, err := balance.NewStorage(config)
	if err != nil {
		log.Fatalf("balance.NewStorage error: %s", err.Error())
	}

	// Create new fiber application service
	service := balance.NewService(config, storage)

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
