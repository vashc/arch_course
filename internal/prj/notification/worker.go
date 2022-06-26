package notification

import (
	"arch_course/internal/prj"
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func NewWorker(config *Config, storage *Storage, client *prj.RabbitClient) *Worker {
	return &Worker{
		config:  config,
		storage: storage,
		client:  client,
	}
}

func (w *Worker) Process(ctx context.Context, queueName string) error {
	queue, err := w.client.Listen(queueName)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-queue:
			err = w.processOne(ctx, msg)
			if err != nil {
				log.Printf("notification processOne error: %s\n", err.Error())
			}
		}
	}
}

func (w *Worker) processOne(ctx context.Context, msg amqp.Delivery) (err error) {
	notification := new(prj.Notification)
	err = json.Unmarshal(msg.Body, &notification)
	if err != nil {
		return err
	}

	return w.storage.UpdateNotificationStatus(notification.ID, StatusCompleted)
}
