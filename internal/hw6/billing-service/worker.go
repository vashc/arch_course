package billing_service

import (
	"arch_course/internal/hw6"
	"encoding/json"
	"errors"
	"fmt"
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
	// Orders queue instantiating
	err := client.CreateQueue(hw6.QueueOrders)
	if err != nil {
		return nil, err
	}

	// Notifications queue instantiating
	if err = client.CreateQueue(hw6.QueueNotifications); err != nil {
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
		log.Printf("listen %s: %s", hw6.QueueOrders, err.Error())
		return err
	}

	go func() {
		delivery := <-msgs
		order := hw6.OrderMessage{}
		err := json.Unmarshal(delivery.Body, &order)
		if err != nil {
			log.Printf("unmarshal OrderMessage: %s", err.Error())
		}

		if err = w.ProcessOrder(&order); err != nil {
			log.Printf("processing OrderMessage: %s", err.Error())
		}
	}()

	<-c
	log.Println("Gracefully shutting down")

	return nil
}

func (w *Worker) ProcessOrder(order *hw6.OrderMessage) error {
	account, err := w.Storage.GetAccountByID(order.UserID)
	if err != nil {
		return fmt.Errorf("%w: %s", hw6.ErrAccountNotFound, err.Error())
	}

	// Compare and update account balance if possible
	err = w.Storage.CompareAndUpdateAccountBalance(account, order.Price)
	if err != nil {
		message := fmt.Sprintf("%s (price: %f)", hw6.NotificationFailInternal, order.Price)

		if errors.Is(err, hw6.ErrInsufficientFunds) {
			message = fmt.Sprintf("%s (price: %f)", hw6.NotificationFailFunds, order.Price)
		}

		err = w.sendNotification(message, account.Email)
		if err != nil {
			return err
		}
	}

	message := fmt.Sprintf("%s (price: %f)", hw6.NotificationSuccess, order.Price)
	return w.sendNotification(message, account.Email)
}

func (w *Worker) sendNotification(message, email string) error {
	notification := &hw6.NotificationMessage{
		Email:   email,
		Payload: message,
	}

	// Publish notification for processing
	return w.Client.Publish(hw6.QueueNotifications, notification)
}
