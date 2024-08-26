package main

import (
	"fmt"

	"github.com/TGPrado/GoScaffoldApi/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Error starting new config.")
	}

	printConfig(cfg)
}

func printConfig(cfg *config.Config) {
	fmt.Printf("App Name: %s\n", cfg.App.Name)
	fmt.Printf("App Version: %s\n", cfg.App.Version)
	fmt.Printf("HTTP Port: %s\n", cfg.HTTP.Port)
	fmt.Printf("Log Level: %s\n", cfg.Log.Level)
	fmt.Printf("DB PoolMax: %d\n", cfg.DB.PoolMax)
	fmt.Printf("DB URL: %s\n", cfg.DB.DBUrl)
	fmt.Printf("DB Timezone: %s\n", cfg.DB.Timezone)
}
