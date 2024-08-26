package database

import (
	"fmt"

	"github.com/TGPrado/GoScaffoldApi/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.DBName,
		cfg.DB.DBPort,
		cfg.DB.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic().Err(err).Msg("Error initializing database")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	return db, nil
}
