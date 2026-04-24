package broker

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewPublisher(amqpURL string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("amqp.Dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("amqp channel: %w", err)
	}
	return &Publisher{conn: conn, ch: ch}, nil
}

func (p *Publisher) Publish(ctx context.Context, queue string, msg any) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if _, err := p.ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	return p.ch.PublishWithContext(ctx, ExchangeDirect, queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (p *Publisher) Close() {
	p.ch.Close()
	p.conn.Close()
}
