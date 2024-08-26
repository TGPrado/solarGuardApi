package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	DB struct {
		PoolMax  int
		Timezone string
		Username string
		Password string
		DBName   string
		DBPort   int
		SSLMode  string
	}
)

func NewConfig() (*Config, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Error opening config file.")
	}

	log.Info().Msg("Config file loaded successfully.")

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panic().Err(err).Msg("Unparsing config file.")
	}

	return &config, nil
}
