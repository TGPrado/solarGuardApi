package usecase

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	ent "github.com/TGPrado/GuardIA/internal/entities"
	"github.com/TGPrado/GuardIA/internal/models"
	repo "github.com/TGPrado/GuardIA/internal/repositories"
	solarz "github.com/TGPrado/GuardIA/pkg/solarZ"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UserUseCase interface {
	ValidateInputs(req ent.UserCreateRequest) error
	ValidateBrandType(req ent.UserCreateRequest) (error, bool)
	CreateUserDatabase(req ent.UserCreateRequest, solarZ int64, brandType bool) error
	CheckIfUserExists(email string) (models.User, error)
	UncreateUser(req ent.UserCreateRequest, automaticBrand bool, user models.User) ent.UncreatedUserResponse
	CreateUser(req ent.UserCreateRequest, automaticBrand bool, user models.User) ent.UncreatedUserResponse
	GetPanelData(idPlant string) ([]solarz.Content, error)
	UpdateUser(user models.User) error
	UpdateUserWithSubscription(end time.Time, email string, planName string, subsId string) error
}

type userUseCase struct {
	deps *deps.Dependencies
	repo repo.UserRepository
}

func NewUserUseCase(deps *deps.Dependencies) UserUseCase {
	repo := repo.NewUserRepository(deps)
	return &userUseCase{deps: deps, repo: repo}
}

func (us userUseCase) ValidateBrandType(req ent.UserCreateRequest) (error, bool) {
	var automatedBrands = map[string]bool{
		"Sunny Postal":    true,
		"Aurora":          true,
		"SolarView":       true,
		"Hoymiles":        true,
		"PV Solar Portal": true,
		"Nexen":           true,
		"Esolar Portal":   true,
		"BYD":             true,
	}

	_, exists := automatedBrands[req.Brand]
	if exists && req.UserInverter != "" && req.PassInverter != "" {
		return nil, true
	}

	if exists {
		return errors.New(
			"O usuário e senha para seu inversor deve ser diferente de nulo",
		), true
	}

	var manualBrands = map[string]bool{
		"Fronius":      true,
		"Fusion Solar": true,
		"Solarman":     true,
		"Growatt":      true,
		"SolarEdge":    true,
		"APsystems":    true,
		"Sunweg":       true,
		"CSI Cloud":    true,
		"Elekeeper":    true,
		"NEP":          true,
		"Sices":        true,
		"Sungrow":      true,
	}

	_, exists = manualBrands[req.Brand]
	if exists {
		return nil, false
	}

	return errors.New("Marca não suportada."), false
}

func (us userUseCase) CreateUserDatabase(req ent.UserCreateRequest, solarZ int64, brandType bool) error {
	return us.repo.Create(req, solarZ, brandType)
}

func (us userUseCase) ValidateInputs(req ent.UserCreateRequest) error {
	regexPhone := `^\(\d{2}\) \d{5}-\d{4}$`
	re := regexp.MustCompile(regexPhone)
	if !re.MatchString(req.Phone) {
		return errors.New("Por favor, forneça um telefone válido.")
	}

	title := cases.Title(language.BrazilianPortuguese)
	req.City = title.String(req.City)
	req.Name = title.String(req.Name)

	if req.Plan != "Plano Basic" && req.Plan != "Plano Premium" {
		return errors.New("Por favor, forneça um plano válido.")
	}

	regexEmail := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re = regexp.MustCompile(regexEmail)
	if !re.MatchString(req.Email) {
		return errors.New("Por favor, forneça um email válido.")
	}

	return nil
}

func (us userUseCase) CheckIfUserExists(email string) (models.User, error) {
	user, err := us.repo.GetByEmail(email)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (us userUseCase) CreateAllUserData(req ent.UserCreateRequest, automaticBrand bool, user models.User) ent.UncreatedUserResponse {
	solarZId, err := solarz.RegisterPanel(req, us.deps.SolarZ)
	if err != nil {
		return ent.UncreatedUserResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
	}
	_ = us.CreateUserDatabase(req, solarZId, automaticBrand)

	return ent.UncreatedUserResponse{
		StatusCode: http.StatusCreated,
		Message:    "Cadastro realizado, você será redirecionado para escolher qual seu painel solar.",
		Id:         solarZId,
	}
}

func (us userUseCase) UncreateUser(req ent.UserCreateRequest, automaticBrand bool, user models.User) ent.UncreatedUserResponse {
	if !automaticBrand || req.Plan == "Plano Premium" {
		// manda msg pro discord
		return ent.UncreatedUserResponse{
			StatusCode: http.StatusOK,
			Message:    "Cadastro realizado, por favor, aguarde um de nossos especialistas entrar em contato",
		}
	}

	return us.CreateAllUserData(req, automaticBrand, user)

}

func (us userUseCase) CreateUser(req ent.UserCreateRequest, automaticBrand bool, user models.User) ent.UncreatedUserResponse {
	if !automaticBrand || user.Plan.Name == "Plano Premium" {
		msg := "Seu cadastro já foi solicitado, entre em contato com o suporte para mais informações."
		return ent.UncreatedUserResponse{
			StatusCode: http.StatusConflict,
			Message:    msg,
		}
	}

	return us.CreateAllUserData(req, automaticBrand, user)
}

func (us userUseCase) GetPanelData(idPlant string) ([]solarz.Content, error) {
	return solarz.GetPlants(idPlant, us.deps.SolarZ)
}

func (us userUseCase) UpdateUser(user models.User) error {
	return us.repo.Update(user)
}

func (us userUseCase) UpdateUserWithSubscription(end time.Time, email string, planName string, subsId string) error {
	model, err := us.repo.GetByEmail(email)
	if err != nil {
		return err
	}

	model.Plan.Name = planName
	model.Plan.RenovationDate = end.Format("02/01/2006 15:04:05")
	model.Plan.Id = subsId

	err = us.repo.Update(model)
	if err != nil {
		return err
	}

	return nil
}
