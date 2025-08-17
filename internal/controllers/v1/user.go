package controllerV1

import (
	"fmt"
	"net/http"

	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	ent "github.com/TGPrado/GuardIA/internal/entities"
	usecase "github.com/TGPrado/GuardIA/internal/useCase"
	"github.com/TGPrado/GuardIA/pkg/helpers"
	solarz "github.com/TGPrado/GuardIA/pkg/solarZ"
	stripelib "github.com/TGPrado/GuardIA/pkg/stripeLib"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Create(c *gin.Context)
	GetPlant(c *gin.Context)
	CreatePlant(c *gin.Context)
}

type userController struct {
	deps    *deps.Dependencies
	useCase usecase.UserUseCase
}

func NewUserController(deps *deps.Dependencies) UserController {
	useCase := usecase.NewUserUseCase(deps)
	return &userController{deps: deps, useCase: useCase}
}

func (us *userController) Create(c *gin.Context) {
	var req ent.UserCreateRequest
	msg, ok := helpers.ValidateInput(&req, c, us.deps)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": msg})
		return
	}

	err, automaticBrand := us.useCase.ValidateBrandType(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	err = us.useCase.ValidateInputs(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	user, err := us.useCase.CheckIfUserExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if user.Email == "" {
		res := us.useCase.UncreateUser(req, automaticBrand, user)
		if res.StatusCode == 201 {
			stripelib.CreateUser(req.Email, req.Phone)
			c.SetCookie("session", req.Email, 3600, "/", "localhost", false, true)
			c.JSON(res.StatusCode, gin.H{"message": res.Message, "id": res.Id})
			return
		}
		if res.StatusCode != 0 {
			c.JSON(res.StatusCode, gin.H{"message": res.Message})
			return
		}
	}

	if user.Plan.Id != "" || user.PanelId != 0 {
		msg := "Por favor, entre em contato com o suporte por um dos canais de comunicação."
		c.JSON(http.StatusConflict, gin.H{"message": msg})
		return
	}

	res := us.useCase.CreateUser(req, automaticBrand, user)
	if res.StatusCode == 201 {
		stripelib.CreateUser(req.Email, req.Phone)
		c.JSON(res.StatusCode, gin.H{"message": res.Message, "id": res.Id})
		return
	}

	c.JSON(res.StatusCode, gin.H{"message": res.Message})
	return
}

func (us *userController) GetPlant(c *gin.Context) {
	idPlant := c.Param("id")
	cookie, err := c.Cookie("session")
	fmt.Println(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Cookie not found"})
		return
	}

	user, err := us.useCase.CheckIfUserExists(cookie)
	if err != nil || user.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if user.SolarzId != 0 {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Você não tem permissão para alterar esse panelId, por favor, entre em contato com o suporte"},
		)
		return
	}

	content, err := us.useCase.GetPanelData(idPlant)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": "Erro encontrado, tente novamente mais tarde"},
		)
		return
	}

	c.JSON(http.StatusOK, content)
}

func (us *userController) CreatePlant(c *gin.Context) {
	id := c.Param("id")
	cookie, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Cookie not found"})
		return
	}

	var req ent.UserCreatePlantRequest
	msg, ok := helpers.ValidateInput(&req, c, us.deps)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": msg})
		return
	}

	user, err := us.useCase.CheckIfUserExists(cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if user.SolarzId != 0 || user.PanelId == 0 {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"message": "Você não tem permissão para alterar esse panelId, por favor, entre em contato com o suporte"},
		)
		return
	}

	err = solarz.CreatePlant(id, req.PlantId, us.deps.SolarZ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro encontrado, tente novamente mais tarde."})
		return
	}

	user.SolarzId = req.PlantId

	err = us.useCase.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro encontrado, tente novamente mais tarde."})
		return
	}

	session, err := stripelib.CreateSession(cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro encontrado, tente novamente mais tarde."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}
