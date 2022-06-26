package hw6

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

func NewClient(config *Config) (*Client, error) {
	connURL := fmt.Sprintf(
		"%s://%s:%s@%s:%s/",
		RabbitProtocol,
		config.RabbitLogin,
		config.RabbitPass,
		config.RabbitHost,
		config.RabbitPort,
	)
	conn, err := amqp.Dial(connURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Client{
		conn,
		ch,
	}, nil
}

func (c *Client) Close() error {
	_ = c.ch.Close()
	return c.conn.Close()
}

func (c *Client) CreateQueue(routingKey string) error {
	// Declaring a queue is idempotent
	_, err := c.ch.QueueDeclare(
		routingKey,
		RabbitDurable,
		RabbitAutoDelete,
		RabbitExclusive,
		RabbitNoWait,
		nil,
	)

	return err
}

func (c *Client) Publish(routingKey string, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.ch.Publish(
		RabbitExchange,
		routingKey,
		RabbitMandatory,
		RabbitImmediate,
		amqp.Publishing{
			ContentType: RabbitContentType,
			Body:        body,
		},
	)
}

func (c *Client) Listen(queueName string) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(
		queueName,
		RabbitConsumer,
		RabbitAutoAck,
		RabbitExclusive,
		RabbitNoLocal,
		RabbitNoWait,
		nil,
	)
}
