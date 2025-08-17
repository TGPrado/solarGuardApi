package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TGPrado/GuardIA/config"
	_ "github.com/TGPrado/GuardIA/docs"
	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	v1 "github.com/TGPrado/GuardIA/internal/routers/v1"
	db "github.com/TGPrado/GuardIA/pkg/database"
	"github.com/TGPrado/GuardIA/pkg/logger"
	httpServer "github.com/TGPrado/GuardIA/pkg/server"
	cfgValidator "github.com/TGPrado/GuardIA/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stripe/stripe-go/v82"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	clientDB, err := db.NewDynamoClient(cfg, l)
	if err != nil {
		fmt.Printf("Não foi possível criar o cliente do DynamoDB: %v\n", err)
		return
	}

	err = db.CreateTableUsers(clientDB, l)
	if err != nil {
		fmt.Printf("Não foi possível criar o cliente do DynamoDB: %v\n", err)
		return
	}

	handler := gin.Default()
	validate, translator := cfgValidator.InitializeValidator(cfg.App.Lang)

	deps := &deps.Dependencies{
		Handler:    handler,
		Logger:     l,
		Validator:  validate,
		Translator: translator,
		DB:         clientDB,
		SolarZ:     cfg.SolarZ,
		Config:     cfg,
	}

	handler.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := "https://guardia.solarguardservices.com.br"
		if origin == "http://localhost:3000" {
			allowOrigin = "http://localhost:3000"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PATCH, DELETE, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	stripe.Key = cfg.Stripe.SecretKey

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1.NewRouter(deps)

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
