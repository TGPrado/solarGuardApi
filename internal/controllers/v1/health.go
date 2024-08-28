package controllerV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

type HealthCheckResponse struct {
	Message string `json:"message"`
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// @Summary Health Check
// @Description Retorna o status de saúde da aplicação
// @Tags Health
// @Produce  json
// @Success 200 {object} HealthCheckResponse
// @Router /health [get]
func (h *HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Ok"})
}
