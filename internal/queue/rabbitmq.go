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
		ch    *amqp.Channel
		queue string
	}
)

func NewRabbitMQPublisher(ch *amqp.Channel, exchange string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		ch:       ch,
		exchange: exchange,
	}
}

func NewRabbitMQSubscriber(ch *amqp.Channel, queue string) *RabbitMQSubscriber {
	return &RabbitMQSubscriber{
		ch:    ch,
		queue: queue,
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

func (m RabbitMQSubscriber) Subscribe(ctx context.Context, handler Handler) error {
	messages, err := m.ch.Consume(
		m.queue,
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
		fmt.Println("message received at", m.queue)
		if err := handler.Handle(ctx, message.Body, message.Headers); err != nil {
			fmt.Println(err)
			if errors.Is(err, event.InvalidPayloadErr) {
				_ = message.Ack(false)
			} else {
				_ = message.Nack(false, true)
			}
		} else {
			_ = message.Ack(false)
		}
	}

	return nil
}
