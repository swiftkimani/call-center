package campaigns

import (
	"context"
	"fmt"

	"github.com/yourorg/callcenter/internal/broker"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	consumer *broker.Consumer
}

func NewWorker(consumer *broker.Consumer) *Worker {
	return &Worker{consumer: consumer}
}

func (w *Worker) Run(ctx context.Context) error {
	log.Info().Msg("campaign worker started")
	return w.consumer.Consume(ctx, broker.QueueCampaignDial, func(d broker.Delivery) error {
		msg, err := broker.UnmarshalMsg[broker.CampaignDialMsg](d)
		if err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		log.Info().
			Str("campaign_id", msg.CampaignID).
			Str("customer_phone", msg.CustomerPhone).
			Msg("campaign dial")
		// Actual dial logic via telephony adapter goes here (wired in main).
		return nil
	})
}
