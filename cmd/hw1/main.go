package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"arch_course/internal/hw1"
)

func main() {
	// Create new fiber application service
	service := hw1.NewService()

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
