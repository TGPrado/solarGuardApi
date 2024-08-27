package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TGPrado/GoScaffoldApi/config"
	"github.com/TGPrado/GoScaffoldApi/pkg/database"
	"github.com/TGPrado/GoScaffoldApi/pkg/logger"
	httpServer "github.com/TGPrado/GoScaffoldApi/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	_, err := database.InitDB(cfg)
	if err != nil {
		l.Panic().Err(err).Msg("Error initializing database")
	}

	handler := gin.Default()
	httpServer := httpServer.New(handler, httpServer.Port(cfg.HTTP.Port))

	if err := gracefulShutdown(httpServer, l); err != nil {
		l.Panic().Err(err).Msg("Error during shutdown")
	}
}

func gracefulShutdown(s *httpServer.Server, l *zerolog.Logger) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info().Msg("Received signal: " + s.String())
	case err := <-s.Notify():
		l.Error().Err(err).Msg("HTTP server error")
	}

	return s.Shutdown()
}
