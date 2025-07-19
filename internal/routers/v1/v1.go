package v1

import (
	controllerV1 "github.com/TGPrado/GuardIA/internal/controllers/v1"
	deps "github.com/TGPrado/GuardIA/internal/dependencies"
)

func NewRouter(deps *deps.Dependencies) {
	healthController := controllerV1.NewHealthController()
	userController := controllerV1.NewUserController(deps)

	deps.Handler.GET("/health", healthController.HealthCheck)
	deps.Handler.POST("/user", userController.Create)
}
