package v1

import (
	controllerV1 "github.com/TGPrado/GoScaffoldApi/internal/controllers/v1"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewRouter(r *gin.Engine, l *zerolog.Logger, db *gorm.DB) {
	healthController := controllerV1.NewHealthController()

	r.GET("/health", healthController.HealthCheck)
}
