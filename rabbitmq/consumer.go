package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	URL  string
	Conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConsumer(url string) *Consumer {
	c := &Consumer{URL: url}
	err := c.connection()
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

func (c *Consumer) connection() error {
	var err error
	c.Conn, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("%s: %s", "Consumer Failed to connect to RabbitMQ", err)
	}
	fmt.Printf("%s\n", "Consumer Successfully connect to RabbitMQ")

	err = c.channel()
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) channel() error {
	var err error
	c.ch, err = c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("%s: %s", "Consumer Failed to open a channel", err)
	}
	fmt.Printf("%s\n", "Consumer Successfully open a channel")
	return nil
}

func (c *Consumer) GetMessage() (<-chan amqp.Delivery, error) {

	if c.Conn.IsClosed() {
		err := c.connection()
		if err != nil {
			return nil, err
		}
	}

	if c.ch.IsClosed() {
		err := c.channel()
		if err != nil {
			return nil, err
		}
	}

	q, err := c.ch.QueueDeclare(
		"queue", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", "Failed to declare a queue", err)
	}

	message, err := c.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", "Failed to register a consumer", err)
	}

	return message, nil
}
