package controllerV1

import (
	"net/http"

	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	ent "github.com/TGPrado/GuardIA/internal/entities"
	"github.com/TGPrado/GuardIA/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Create(c *gin.Context)
}

type userController struct {
	deps *deps.Dependencies
}

func NewUserController(deps *deps.Dependencies) UserController {
	return &userController{deps: deps}
}

func (us *userController) Create(c *gin.Context) {
	var req ent.UserCreateRequest
	msg, ok := helpers.ValidateInput(&req, c, us.deps)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": msg})
		return
	}
}
