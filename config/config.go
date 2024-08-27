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
		Host     string
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

	viper.BindEnv("database.poolMax", "POOL_MAX")
	viper.BindEnv("database.username", "DB_USERNAME")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbName", "DB_NAME")
	viper.BindEnv("database.dbPort", "DB_PORT")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.poolsslModeMax", "DB_SSL")
	viper.BindEnv("database.dbTimezone", "DB_TIMEZONE")

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
