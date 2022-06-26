package main

import (
	"arch_course/internal/prj"
	"arch_course/internal/prj/wallet"
	"context"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create config, e.g. via environment variables
	config, err := wallet.NewConfig()
	if err != nil {
		log.Fatalf("wallet.NewConfig error: %s", err.Error())
	}

	// Initialize storage
	storage, err := wallet.NewStorage(config)
	if err != nil {
		log.Fatalf("wallet.NewStorage error: %s", err.Error())
	}

	// Initialize Rabbit client and connection
	client, err := prj.NewClient(config.RabbitConfig)
	if err != nil {
		log.Fatalf("prj.NewClient error: %s", err.Error())
	}
	defer client.Close()

	// Instantiate internal queue
	err = client.CreateQueue(prj.QueueSagaSteps)
	if err != nil {
		log.Fatalf("client.CreateQueue error: %s", err.Error())
	}

	// Create new fiber application service
	service := wallet.NewService(config, storage, client)

	// Instantiate routes
	service.InstantiateRoutes()

	// Start saga worker
	ctx, cancel := context.WithCancel(context.Background())
	worker := wallet.NewWorker(config, storage, client)
	go func() {
		err = worker.Process(ctx, prj.QueueSagaSteps)
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
	if err = service.Start(":8000"); err != nil {
		log.Panicf("service.Start error: %s", err.Error())
	}

	log.Println("Running cleanup")
}
