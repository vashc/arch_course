package bcgateway

import (
	"arch_course/internal/prj"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)

func NewWorker(config *Config, storage *Storage, rabbitClient *prj.RabbitClient) *Worker {
	return &Worker{
		config:  config,
		storage: storage,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		rabbitClient: rabbitClient,
	}
}

func (w *Worker) Process(ctx context.Context, queueName string) error {
	queue, err := w.rabbitClient.Listen(queueName)
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
				log.Printf("bcgateway processOne error: %s\n", err.Error())
			}
		}
	}
}

func (w *Worker) processOne(ctx context.Context, msg amqp.Delivery) (err error) {
	order := new(prj.BcgatewayOrder)
	err = json.Unmarshal(msg.Body, &order)
	if err != nil {
		return err
	}

	time.Sleep(7 * time.Second)

	err = w.updateBalances(order)
	if err != nil {
		return err
	}

	if !order.Compensate {
		err = w.rabbitClient.Publish(
			prj.QueueSagaSteps,
			prj.SagaStep{
				OrderID: order.OrderID,
				Type:    prj.SagaTypeBcgateway,
				Status:  StatusCompleted,
			},
		)
		if err != nil {
			return err
		}
	}

	return w.storage.UpdateOrderStatus(order.ID, StatusCompleted)
}

func (w *Worker) updateBalance(userID int64, cryptoAmount, fiatAmount float64) (err error) {
	// Decrease user fiat balance
	balanceURL := fmt.Sprintf(
		"http://%s:%d/wallet/%d",
		w.config.BalanceHost,
		w.config.BalancePort,
		userID,
	)

	balanceResp, err := prj.DoRequest(
		w.client,
		balanceURL,
		http.MethodPatch,
		&prj.Wallet{
			CryptoAmount: cryptoAmount,
			FiatAmount:   fiatAmount,
		},
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		log.Printf("balance DoRequest error: %s\n", err.Error())
		code := http.StatusInternalServerError
		return errors.New(http.StatusText(code))
	}

	if balanceResp.StatusCode != http.StatusOK {
		log.Printf("balanceResp.StatusCode: %d\n", balanceResp.StatusCode)
		code := balanceResp.StatusCode
		return errors.New(http.StatusText(code))
	}

	return nil
}

func (w *Worker) updateBalances(order *prj.BcgatewayOrder) (err error) {
	// Decrease hot wallet crypto amount
	err = w.updateBalance(prj.HotWalletUserID, -prj.Btof(order.Compensate)*order.CryptoAmount, 0)
	if err != nil {
		return err
	}

	// Increase user crypto amount
	return w.updateBalance(order.AcquirerUserID, prj.Btof(order.Compensate)*order.CryptoAmount, 0)
}
