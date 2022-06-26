package wallet

import (
	"arch_course/internal/prj"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
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
				log.Printf("wallet processOne error: %s\n", err.Error())
			}
		}
	}
}

func (w *Worker) processOne(ctx context.Context, msg amqp.Delivery) (err error) {
	step := new(prj.SagaStep)
	err = json.Unmarshal(msg.Body, &step)
	if err != nil {
		return err
	}

	// A step has been completed
	var stepsRemain int
	stepsRemain, err = w.storage.DecreaseCurrentStepNumber(step.OrderID)
	if err != nil {
		return err
	}

	// A step has been failed
	if step.Status == StatusFailed {
		err = w.storage.FailOrderProcessingStep(step.OrderID, step.Type)
		if err != nil {
			return err
		}
	}

	if stepsRemain == 0 {
		var failedSteps map[uint8]struct{}

		failedSteps, err = w.storage.GetFailedOrderProcessingSteps(step.OrderID)
		if err != nil {
			return err
		}

		status := StatusCompleted
		if len(failedSteps) > 0 {
			// If we have some completed orders, we should compensate them
			status = StatusFailed

			var order *Order
			order, err = w.storage.GetOrderByID(step.OrderID)
			if err != nil {
				return err
			}

			err = w.createCompensatingOrders(order, failedSteps)
			if err != nil {
				return err
			}
		}

		err = w.formAndSendNotification(step.OrderID, status)
		if err != nil {
			return err
		}

		return w.storage.UpdateOrderStatus(step.OrderID, status)
	}

	return nil
}

func (w *Worker) getUserData(userID int64) (user *prj.User, err error) {
	authURL := fmt.Sprintf(
		"http://%s:%d/user/%d",
		w.config.AuthHost,
		w.config.AuthPort,
		userID,
	)

	authResp, err := prj.DoRequest(w.client, authURL, http.MethodGet, nil, nil)
	if err != nil {
		log.Printf("auth DoRequest error: %s\n", err.Error())
		return nil, err
	}
	defer authResp.Body.Close()

	if authResp.StatusCode != http.StatusOK {
		log.Printf("authResp.StatusCode: %d\n", authResp.StatusCode)
		return nil, errors.New(http.StatusText(authResp.StatusCode))
	}

	authBody, err := ioutil.ReadAll(authResp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %s\n", err.Error())
		return nil, err
	}

	err = json.Unmarshal(authBody, &user)
	if err != nil {
		log.Printf("auth user (%s) BodyParser error: %s", string(authBody), err.Error())
		return nil, err
	}

	return user, nil
}

func (w *Worker) sendNotification(payload, email string, orderID int64) (code int, err error) {
	// Make request to notification service
	notificationURL := fmt.Sprintf(
		"http://%s:%d/send/%d",
		w.config.NotificationHost,
		w.config.NotificationPort,
		orderID,
	)

	notification := &prj.Notification{
		Email:   email,
		Payload: payload,
	}

	notificationResp, err := prj.DoRequest(
		w.client,
		notificationURL,
		http.MethodPost,
		notification,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		log.Printf("notification DoRequest error: %s\n", err.Error())
		return http.StatusInternalServerError, err
	}
	defer notificationResp.Body.Close()

	if notificationResp.StatusCode != http.StatusOK {
		log.Printf("notificationResp.StatusCode: %d\n", notificationResp.StatusCode)
		code = notificationResp.StatusCode
		return code, nil
	}

	return http.StatusOK, nil
}

func (w *Worker) formAndSendNotification(orderID int64, status string) (err error) {
	// Get order data
	var order *Order
	order, err = w.storage.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	// Get user data
	var user *prj.User
	user, err = w.getUserData(order.UserID)
	if err != nil {
		return err
	}

	// Form notification payload
	payload := fmt.Sprintf(
		"%s Order data: {fiat amount: %f, crypto amount: %f}",
		NotificationOrderCompleted,
		order.CryptoAmount,
		order.FiatAmount,
	)
	if status == StatusFailed {
		payload = fmt.Sprintf(
			"%s Order data: {fiat amount: %f, crypto amount: %f}",
			NotificationOrderFailed,
			order.CryptoAmount,
			order.FiatAmount,
		)
	}

	// Send notification
	_, err = w.sendNotification(payload, user.Email, orderID)
	return err
}

func (w *Worker) createCompensatingOrders(order *Order, steps map[uint8]struct{}) (err error) {
	// Make compensation orders
	orderSteps := order.GetSteps()
	for step := range orderSteps {
		if _, ok := steps[step]; !ok {
			_, err = orderSteps[step](w, order, true)
			if err != nil {
				return err
			}
		}
	}

	if order.Type == TypeBuyOrder {
		// Release fiat
		_, err = holdFiat(w, order, true)
		return err
	}

	// Release crypto
	_, err = holdCrypto(w, order, true)
	return err
}

func (w *Worker) HttpClient() *http.Client {
	return w.client
}

func (w *Worker) BalanceURI() string {
	return fmt.Sprintf("http://%s:%d", w.config.BalanceHost, w.config.BalancePort)
}

func (w *Worker) BcgatewayURI() string {
	return fmt.Sprintf("http://%s:%d", w.config.BcgatewayHost, w.config.BcgatewayPort)
}

func (w *Worker) ExchangerURI() string {
	return fmt.Sprintf("http://%s:%d", w.config.ExchangerHost, w.config.ExchangerPort)
}
