package config

import (
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type (
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		DB      DB
		SolarZ  SolarZ
		Stripe  Stripe
		Discord Discord
	}

	App struct {
		Name        string
		Version     string
		Lang        string
		Host        string
		Environment string
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	SolarZ struct {
		Email    string
		Password string
	}

	DB struct {
		Region   string
		Endpoint string
	}

	Stripe struct {
		PubKey        string
		SecretKey     string
		WebhookSecret string
	}

	Discord struct {
		Webhook string
	}
)

func NewConfig() (*Config, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AutomaticEnv()

	viper.BindEnv("db.host", "DB_HOST")

	viper.BindEnv("database.region", "DB_REGION")
	viper.BindEnv("database.endpoint", "DB_ENDPOINT")

	viper.BindEnv("solarz.email", "SOLARZ_EMAIL")
	viper.BindEnv("solarz.password", "SOLARZ_PASSWORD")

	viper.BindEnv("stripe.pubKey", "STRIPE_PUB_KEY")
	viper.BindEnv("stripe.secretKey", "STRIPE_SECRET_KEY")
	viper.BindEnv("stripe.webhookSecret", "WEBHOOK_SECRET_STRIPE")

	viper.BindEnv("app.environment", "ENVIRONMENT")

	viper.BindEnv("discord.webhook", "DISCORD_WEBHOOK")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Panic().Err(err).Msg("Error opening config file.")
	}

	logger.Info().Msg("Config file loaded successfully.")

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Panic().Err(err).Msg("Unparsing config file.")
	}

	return &config, nil
}
