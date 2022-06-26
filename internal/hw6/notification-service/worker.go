package notification_service

import (
	"arch_course/internal/hw6"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewWorker(
	config *hw6.Config,
	storage *Storage,
	client *hw6.Client,
) (*Worker, error) {
	// Notifications queue instantiating
	if err := client.CreateQueue(hw6.QueueNotifications); err != nil {
		return nil, err
	}

	return &Worker{
		Config:  config,
		Storage: storage,
		Client:  client,
	}, nil
}

func (w *Worker) Start() error {
	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	msgs, err := w.Client.Listen(hw6.QueueOrders)
	if err != nil {
		log.Printf("listen %s: %s", hw6.QueueNotifications, err.Error())
		return err
	}

	go func() {
		delivery := <-msgs
		notification := hw6.NotificationMessage{}
		err := json.Unmarshal(delivery.Body, &notification)
		if err != nil {
			log.Printf("unmarshal NotificationMessage: %s", err.Error())
		}

		if err = w.ProcessNotification(&notification); err != nil {
			log.Printf("processing NotificationMessage: %s", err.Error())
		}
	}()

	<-c
	log.Println("Gracefully shutting down")

	return nil
}

func (w *Worker) ProcessNotification(order *hw6.NotificationMessage) error {
	message := &Message{
		Email:   order.Email,
		Payload: order.Payload,
	}

	return w.Storage.SendMessage(message)
}
