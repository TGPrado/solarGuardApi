package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TGPrado/GoScaffoldApi/config"
	_ "github.com/TGPrado/GoScaffoldApi/docs"
	v1 "github.com/TGPrado/GoScaffoldApi/internal/routers/v1"
	"github.com/TGPrado/GoScaffoldApi/pkg/database"
	"github.com/TGPrado/GoScaffoldApi/pkg/logger"
	httpServer "github.com/TGPrado/GoScaffoldApi/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	db, err := database.InitDB(cfg)
	if err != nil {
		l.Panic().Err(err).Msg("Error initializing database")
	}

	handler := gin.Default()
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1.NewRouter(handler, l, db)

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
