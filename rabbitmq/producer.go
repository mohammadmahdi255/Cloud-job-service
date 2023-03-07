package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	URL  string
	Conn *amqp.Connection
	ch   *amqp.Channel
}

func NewProducer(url string) *Producer {
	p := &Producer{URL: url}
	err := p.connection()
	if err != nil {
		log.Println(err.Error())
	}
	return p
}

func (p *Producer) connection() error {
	var err error
	p.Conn, err = amqp.Dial(p.URL)
	if err != nil {
		return fmt.Errorf("%s: %s", "Producer Failed to connect to RabbitMQ", err)
	}
	fmt.Printf("%s\n", "Producer Successfully connect to RabbitMQ")

	err = p.channel()
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) channel() error {
	var err error
	p.ch, err = p.Conn.Channel()
	if err != nil {
		return fmt.Errorf("%s: %s", "Producer Failed to open a channel", err)
	}
	fmt.Printf("%s\n", "Producer Successfully open a channel")
	return nil
}

func (p *Producer) Send(body string) error {

	if p.Conn.IsClosed() {
		err := p.connection()
		if err != nil {
			return err
		}
	}

	if p.ch.IsClosed() {
		err := p.channel()
		if err != nil {
			return err
		}
	}

	q, err := p.ch.QueueDeclare(
		"queue", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to declare a queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to publish a message", err)
	}
	log.Printf("Sent %s\n", body)

	return nil
}
