package main

import (
	"github.com/TGPrado/GoScaffoldApi/config"
	"github.com/TGPrado/GoScaffoldApi/internal/app"
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
