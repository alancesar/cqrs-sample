package queue

import (
	"context"
	"cqrs-sample/pkg/event"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Handler interface {
		Handle(ctx context.Context, body []byte, headers map[string]interface{}) error
	}

	RabbitMQPublisher struct {
		ch       *amqp.Channel
		exchange string
	}

	RabbitMQSubscriber struct {
		conn *amqp.Connection
	}
)

func NewRabbitMQPublisher(ch *amqp.Channel, exchange string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		ch:       ch,
		exchange: exchange,
	}
}

func NewRabbitMQSubscriber(conn *amqp.Connection) *RabbitMQSubscriber {
	return &RabbitMQSubscriber{
		conn: conn,
	}
}

func (p RabbitMQPublisher) Publish(ctx context.Context, message event.Message, e event.Event) error {
	err := p.ch.PublishWithContext(ctx,
		p.exchange,
		string(e),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message.Body,
			Headers:     message.Headers,
		})

	fmt.Println("message received at", p.exchange, e)
	return err
}

func (m RabbitMQSubscriber) Subscribe(ctx context.Context, queue string, handler Handler) error {
	channel, err := m.conn.Channel()
	if err != nil {
		return err
	}

	defer func() {
		_ = channel.Close()
	}()

	messages, err := channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for message := range messages {
		fmt.Println("message received at", queue)
		err := handler.Handle(ctx, message.Body, message.Headers)
		if err != nil {
			fmt.Println(err)
		}

		if err == nil || errors.Is(err, event.InvalidPayloadErr) {
			_ = message.Ack(false)
		} else {
			_ = message.Nack(false, true)
		}
	}

	return nil
}
