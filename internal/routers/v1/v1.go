package v1

import (
	controllerV1 "github.com/TGPrado/GuardIA/internal/controllers/v1"
	deps "github.com/TGPrado/GuardIA/internal/dependencies"
)

func NewRouter(deps *deps.Dependencies) {
	healthController := controllerV1.NewHealthController()
	userController := controllerV1.NewUserController(deps)
	paymentsController := controllerV1.NewPaymentsController(deps)

	deps.Handler.GET("/health", healthController.HealthCheck)
	deps.Handler.POST("api/user", userController.Create)
	deps.Handler.GET("api/user/:id", userController.GetPlant)
	deps.Handler.POST("api/user/:id", userController.CreatePlant)

	deps.Handler.POST("api/webhook", paymentsController.Webhook)
}
