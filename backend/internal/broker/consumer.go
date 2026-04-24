package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type Delivery = amqp.Delivery

type Consumer struct {
	amqpURL string
}

func NewConsumer(amqpURL string) *Consumer {
	return &Consumer{amqpURL: amqpURL}
}

// Consume runs a consumer loop with automatic reconnection.
func (c *Consumer) Consume(ctx context.Context, queue string, handler func(Delivery) error) error {
	backoff := time.Second
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := c.consumeOnce(ctx, queue, handler); err != nil {
			log.Error().Err(err).Str("queue", queue).Dur("backoff", backoff).Msg("consumer error, reconnecting")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				if backoff < 60*time.Second {
					backoff *= 2
				}
			}
		} else {
			return nil
		}
	}
}

func (c *Consumer) consumeOnce(ctx context.Context, queue string, handler func(Delivery) error) error {
	conn, err := amqp.Dial(c.amqpURL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	connClosed := conn.NotifyClose(make(chan *amqp.Error, 1))
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-connClosed:
			if err != nil {
				return fmt.Errorf("connection closed: %w", err)
			}
			return nil
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			if err := handler(d); err != nil {
				log.Error().Err(err).Str("queue", queue).Msg("handler error, nacking")
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}
}

func UnmarshalMsg[T any](d Delivery) (T, error) {
	var v T
	if err := json.Unmarshal(d.Body, &v); err != nil {
		return v, fmt.Errorf("unmarshal: %w", err)
	}
	return v, nil
}
