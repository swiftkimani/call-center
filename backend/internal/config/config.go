package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	AMQPURL     string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool

	JWTSecret        string
	JWTRefreshSecret string

	ATAPIKey        string
	ATUsername      string
	ATWebhookSecret string

	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioWebhookSecret string

	TelephonyProvider string // "africas_talking" | "twilio"
	BaseURL           string

	Port             string
	LogLevel         string
	DialingHourStart int
	DialingHourEnd   int
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		RedisURL:         os.Getenv("REDIS_URL"),
		AMQPURL:          os.Getenv("AMQP_URL"),
		MinioEndpoint:    os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey:   os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey:   os.Getenv("MINIO_SECRET_KEY"),
		MinioBucket:      getEnvOrDefault("MINIO_BUCKET", "recordings"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		ATAPIKey:         os.Getenv("AT_API_KEY"),
		ATUsername:       os.Getenv("AT_USERNAME"),
		ATWebhookSecret:  os.Getenv("AT_WEBHOOK_SECRET"),
		TwilioAccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		TelephonyProvider: getEnvOrDefault("TELEPHONY_PROVIDER", "africas_talking"),
		BaseURL:          getEnvOrDefault("BASE_URL", "http://localhost:8080"),
		Port:             getEnvOrDefault("PORT", "8080"),
		LogLevel:         getEnvOrDefault("LOG_LEVEL", "info"),
	}

	var err error
	cfg.MinioUseSSL, err = strconv.ParseBool(getEnvOrDefault("MINIO_USE_SSL", "false"))
	if err != nil {
		return nil, fmt.Errorf("MINIO_USE_SSL: %w", err)
	}

	cfg.DialingHourStart, err = strconv.Atoi(getEnvOrDefault("DIALING_HOUR_START", "8"))
	if err != nil {
		return nil, fmt.Errorf("DIALING_HOUR_START: %w", err)
	}

	cfg.DialingHourEnd, err = strconv.Atoi(getEnvOrDefault("DIALING_HOUR_END", "20"))
	if err != nil {
		return nil, fmt.Errorf("DIALING_HOUR_END: %w", err)
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	var errs []error
	required := map[string]string{
		"DATABASE_URL":       c.DatabaseURL,
		"REDIS_URL":          c.RedisURL,
		"AMQP_URL":           c.AMQPURL,
		"MINIO_ENDPOINT":     c.MinioEndpoint,
		"MINIO_ACCESS_KEY":   c.MinioAccessKey,
		"MINIO_SECRET_KEY":   c.MinioSecretKey,
		"JWT_SECRET":         c.JWTSecret,
		"JWT_REFRESH_SECRET": c.JWTRefreshSecret,
	}
	for k, v := range required {
		if v == "" {
			errs = append(errs, fmt.Errorf("missing required env var: %s", k))
		}
	}
	return errors.Join(errs...)
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
