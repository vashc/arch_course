package hw6

import (
	"github.com/segmentio/kafka-go"
	"github.com/streadway/amqp"
)

type Config struct {
	DBURI    string `envconfig:"DATABASE_URI" required:"true"`
	Hostname string `envconfig:"HOSTNAME"`

	RabbitHost  string `envconfig:"RABBIT_HOST" required:"true"`
	RabbitPort  string `envconfig:"RABBIT_PORT" required:"true"`
	RabbitLogin string `envconfig:"RABBIT_LOGIN" required:"true"`
	RabbitPass  string `envconfig:"RABBIT_PASS" required:"true"`
}

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type Client1 struct {
	*kafka.Reader
	*kafka.Writer
}

type OrderMessage struct {
	UserID int64   `json:"user_id"`
	Price  float64 `json:"price"`
}

type NotificationMessage struct {
	Email   string `json:"email"`
	Payload string `json:"payload"`
}
