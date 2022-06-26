package prj

import "github.com/streadway/amqp"

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Wallet struct {
	UserID       int64   `json:"user_id"`
	CryptoAmount float64 `json:"crypto_amount"`
	FiatAmount   float64 `json:"fiat_amount"`
}

type Notification struct {
	ID      int64  `json:"id"`
	OrderID int64  `json:"order_id"`
	Email   string `json:"email"`
	Payload string `json:"payload"`
	Status  string `json:"status"`
}

type ExchangeOrder struct {
	ID             int64   `json:"id"`
	UUID           string  `json:"uuid"`
	AcquirerUserID int64   `json:"acquirer_user_id"`
	OrderID        int64   `json:"order_id"`
	Type           string  `json:"type"`
	FiatAmount     float64 `json:"fiat_amount"`
	CryptoAmount   float64 `json:"crypto_amount"`
	Compensate     bool    `json:"compensate"`
	Status         string  `json:"status"`
}

type BcgatewayOrder struct {
	ID             int64   `json:"id"`
	UUID           string  `json:"uuid"`
	AcquirerUserID int64   `json:"acquirer_user_id"`
	OrderID        int64   `json:"order_id"`
	CryptoAmount   float64 `json:"crypto_amount"`
	Compensate     bool    `json:"compensate"`
	Status         string  `json:"status"`
}

type SagaStep struct {
	OrderID int64  `json:"order_uuid"`
	Type    int    `json:"step_type"`
	Status  string `json:"status"`
}

type RabbitConfig struct {
	RabbitHost  string `envconfig:"RABBIT_HOST" required:"true"`
	RabbitPort  string `envconfig:"RABBIT_PORT" required:"true"`
	RabbitLogin string `envconfig:"RABBIT_LOGIN" required:"true"`
	RabbitPass  string `envconfig:"RABBIT_PASS" required:"true"`
}

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}
