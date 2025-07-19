package main

import (
	"github.com/TGPrado/GuardIA/config"
	"github.com/TGPrado/GuardIA/internal/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Error starting new config.")
	}

	app.Run(cfg)
}
