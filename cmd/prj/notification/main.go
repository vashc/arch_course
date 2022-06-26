package main

import (
	"arch_course/internal/prj"
	"arch_course/internal/prj/notification"
	"context"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create config, e.g. via environment variables
	config, err := notification.NewConfig()
	if err != nil {
		log.Fatalf("notification.NewConfig error: %s", err.Error())
	}

	// Initialize storage
	storage, err := notification.NewStorage(config)
	if err != nil {
		log.Fatalf("notification.NewStorage error: %s", err.Error())
	}

	// Initialize Rabbit client and connection
	client, err := prj.NewClient(config.RabbitConfig)
	if err != nil {
		log.Fatalf("prj.NewClient error: %s", err.Error())
	}
	defer client.Close()

	// Instantiate internal queue
	err = client.CreateQueue(prj.QueueNotifications)
	if err != nil {
		log.Fatalf("client.CreateQueue error: %s", err.Error())
	}

	// Create new fiber application service
	service := notification.NewService(config, storage, client)

	// Instantiate routes
	service.InstantiateRoutes()

	// Start worker
	ctx, cancel := context.WithCancel(context.Background())
	worker := notification.NewWorker(config, storage, client)
	go func() {
		err = worker.Process(ctx, prj.QueueNotifications)
		if err != nil {
			log.Fatalf("worker.Process error: %s", err.Error())
		}
	}()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
		log.Println("Gracefully shutting down")
	}()

	// Start handling requests
	if err := service.Start(":8000"); err != nil {
		log.Panicf("service.Start error: %s", err.Error())
	}

	log.Println("Running cleanup")
}
