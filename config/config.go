package config

import (
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type (
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		DB   DB
	}

	App struct {
		Name    string
		Version string
		Lang    string
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	DB struct {
		Region   string
		Endpoint string
	}
)

func NewConfig() (*Config, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AutomaticEnv()

	viper.BindEnv("database.region", "DB_REGION")
	viper.BindEnv("database.endpoint", "DB_ENDPOINT")

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
