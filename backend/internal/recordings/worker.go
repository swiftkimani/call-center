package recordings

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/callcenter/internal/broker"
	"github.com/yourorg/callcenter/internal/calls"
	"github.com/yourorg/callcenter/internal/storage"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	consumer *broker.Consumer
	store    *storage.Client
	callsSvc *calls.Service
}

func NewWorker(consumer *broker.Consumer, store *storage.Client, callsSvc *calls.Service) *Worker {
	return &Worker{consumer: consumer, store: store, callsSvc: callsSvc}
}

func (w *Worker) Run(ctx context.Context) error {
	log.Info().Msg("recordings worker started")
	return w.consumer.Consume(ctx, broker.QueueRecordingReady, func(d broker.Delivery) error {
		msg, err := broker.UnmarshalMsg[broker.RecordingReadyMsg](d)
		if err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		return w.process(ctx, msg)
	})
}

func (w *Worker) process(ctx context.Context, msg broker.RecordingReadyMsg) error {
	callID, err := uuid.Parse(msg.CallID)
	if err != nil {
		return fmt.Errorf("invalid call_id: %w", err)
	}

	resp, err := http.Get(msg.ProviderRecordingURL) //nolint:noctx
	if err != nil {
		return fmt.Errorf("download recording: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	now := time.Now().UTC()
	key := fmt.Sprintf("recordings/%d/%02d/%02d/%s.opus", now.Year(), now.Month(), now.Day(), callID)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if err := w.store.Upload(ctx, key, strings.NewReader(string(body)), int64(len(body)), "audio/ogg"); err != nil {
		return fmt.Errorf("upload: %w", err)
	}

	if err := w.callsSvc.UpdateRecordingURL(ctx, callID, key); err != nil {
		return fmt.Errorf("UpdateRecordingURL: %w", err)
	}

	log.Info().Str("call_id", msg.CallID).Str("key", key).Msg("recording stored")
	return nil
}
