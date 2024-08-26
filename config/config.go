package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type (
	Config struct {
		App  App  `json:"app"`
		HTTP HTTP `json:"http"`
		Log  Log  `json:"logger"`
		DB   DB   `json:"database"`
	}

	App struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	HTTP struct {
		Port string `json:"port"`
	}

	Log struct {
		Level string `json:"level"`
	}

	DB struct {
		PoolMax  int    `json:"poolMax"`
		DBUrl    string `json:"dbUrl"`
		Timezone string `json:"dbTimezone"`
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
