package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yourorg/callcenter/internal/agents"
	"github.com/yourorg/callcenter/internal/broker"
	"github.com/yourorg/callcenter/internal/calls"
	"github.com/yourorg/callcenter/internal/campaigns"
	"github.com/yourorg/callcenter/internal/config"
	"github.com/yourorg/callcenter/internal/customers"
	"github.com/yourorg/callcenter/internal/db"
	"github.com/yourorg/callcenter/internal/hub"
	redisclient "github.com/yourorg/callcenter/internal/redis"
	"github.com/yourorg/callcenter/internal/queues"
	"github.com/yourorg/callcenter/internal/reports"
	"github.com/yourorg/callcenter/internal/storage"
	"github.com/yourorg/callcenter/internal/telephony/africas_talking"
	"github.com/yourorg/callcenter/internal/telephony/twilio"
	httpTransport "github.com/yourorg/callcenter/internal/transport/http"
	"github.com/yourorg/callcenter/internal/transport/http/handlers"
	"github.com/yourorg/callcenter/internal/transport/http/webhooks"
)

func main() {
	// Logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	if cfg.LogLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Database
	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect db")
	}
	defer pool.Close()

	if err := db.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("run migrations")
	}
	log.Info().Msg("migrations: up to date")

	// Redis
	rdb, err := redisclient.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connect redis")
	}
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("redis ping")
	}

	// Message broker
	publisher, err := broker.NewPublisher(cfg.AMQPURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect rabbitmq publisher")
	}
	defer publisher.Close()

	// Object storage
	store, err := storage.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connect minio")
	}
	if err := store.EnsureBucket(ctx); err != nil {
		log.Warn().Err(err).Msg("ensure bucket (non-fatal)")
	}

	// Telephony adapter
	var telephonyAdapter interface {
		// reuse the concrete type, assign to interface below
	}
	_ = telephonyAdapter

	var atAdapter = africas_talking.NewAdapter(cfg.ATAPIKey, cfg.ATUsername, cfg.ATWebhookSecret)
	var twilioAdapter = twilio.NewAdapter(cfg.TwilioAccountSID, cfg.TwilioAuthToken, cfg.TwilioWebhookSecret)
	_ = twilioAdapter

	// Domain services
	agentsSvc := agents.NewService(pool, rdb)
	queueSvc := queues.NewService(pool, rdb)
	queueRouter := queues.NewRouter(pool, rdb)
	customersSvc := customers.NewService(pool)
	callsSvc := calls.NewService(pool, rdb, atAdapter, cfg.DialingHourStart, cfg.DialingHourEnd)
	campaignsSvc := campaigns.NewService(pool)
	reportsSvc := reports.NewService(pool)

	// WebSocket hub
	wsHub := hub.New(rdb)
	go wsHub.Run(ctx)

	// HTTP handlers
	authHandler := handlers.NewAuthHandler(cfg, pool)
	meHandler := handlers.NewMeHandler(pool)
	agentsHandler := handlers.NewAgentsHandler(agentsSvc)
	campaignsHandler := handlers.NewCampaignsHandler(campaignsSvc)
	callsHandler := handlers.NewCallsHandler(callsSvc, customersSvc, store)
	customersHandler := handlers.NewCustomersHandler(customersSvc)
	reportsHandler := handlers.NewReportsHandler(reportsSvc)
	supervisorHandler := handlers.NewSupervisorHandler(wsHub, queueSvc)

	atWebhookHandler := webhooks.NewATHandler(
		atAdapter, cfg, callsSvc, customersSvc,
		agentsSvc, queueSvc, queueRouter, wsHub, publisher,
	)

	router := httpTransport.NewRouter(httpTransport.Deps{
		Cfg:               cfg,
		AuthHandler:       authHandler,
		MeHandler:         meHandler,
		AgentsHandler:     agentsHandler,
		CampaignsHandler:  campaignsHandler,
		CallsHandler:      callsHandler,
		CustomersHandler:  customersHandler,
		ReportsHandler:    reportsHandler,
		SupervisorHandler: supervisorHandler,
		ATWebhookHandler:  atWebhookHandler,
		Hub:               wsHub,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-quit
	log.Info().Msg("shutting down...")
	cancel()

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
	}
	log.Info().Msg("stopped")
}
