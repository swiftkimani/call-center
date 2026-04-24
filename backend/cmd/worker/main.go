package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yourorg/callcenter/internal/broker"
	"github.com/yourorg/callcenter/internal/calls"
	"github.com/yourorg/callcenter/internal/campaigns"
	"github.com/yourorg/callcenter/internal/config"
	"github.com/yourorg/callcenter/internal/db"
	redisclient "github.com/yourorg/callcenter/internal/redis"
	"github.com/yourorg/callcenter/internal/recordings"
	"github.com/yourorg/callcenter/internal/storage"
	"github.com/yourorg/callcenter/internal/telephony/africas_talking"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect db")
	}
	defer pool.Close()

	rdb, err := redisclient.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connect redis")
	}

	store, err := storage.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connect minio")
	}
	store.EnsureBucket(ctx)

	atAdapter := africas_talking.NewAdapter(cfg.ATAPIKey, cfg.ATUsername, cfg.ATWebhookSecret)
	callsSvc := calls.NewService(pool, rdb, atAdapter, cfg.DialingHourStart, cfg.DialingHourEnd)

	consumer := broker.NewConsumer(cfg.AMQPURL)

	recWorker := recordings.NewWorker(consumer, store, callsSvc)
	campWorker := campaigns.NewWorker(consumer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := recWorker.Run(ctx); err != nil {
			log.Error().Err(err).Msg("recordings worker stopped")
		}
	}()

	go func() {
		defer wg.Done()
		if err := campWorker.Run(ctx); err != nil {
			log.Error().Err(err).Msg("campaign worker stopped")
		}
	}()

	<-quit
	log.Info().Msg("worker shutting down...")
	cancel()
	wg.Wait()
	log.Info().Msg("worker stopped")
}
