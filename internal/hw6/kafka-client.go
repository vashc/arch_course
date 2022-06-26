package hw6

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func NewClient1(config *Config, topic string) *Client {
	reader := newReader(config, topic)
	writer := newWriter(config, topic)

	return &Client{Reader: reader, Writer: writer}
}

func newReader(cfg *Config, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)},
		Topic:     topic,
		Partition: KafkaPartition,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
}

func newWriter(cfg *Config, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (c *Client) CloseReader() error {
	return c.Reader.Close()
}

func (c *Client) CloseWriter() error {
	return c.Writer.Close()
}

func (c *Client) Read(input chan []byte) (err error) {
	var message kafka.Message

	for {
		message, err = c.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		input <- message.Value
	}
}

func (c *Client) Write(output chan []byte) (err error) {
	for {
		rawMessage := <-output

		err = c.WriteMessages(
			context.Background(),
			kafka.Message{
				Value: rawMessage,
			},
		)
		if err != nil {
			return err
		}
	}
}
