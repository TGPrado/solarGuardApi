package repositories

import (
	"context"
	"fmt"

	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	ent "github.com/TGPrado/GuardIA/internal/entities"
	"github.com/TGPrado/GuardIA/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepository interface {
	GetByEmail(email string) (models.User, error)
	Create(req ent.UserCreateRequest, solarZ int64, brandType bool) error
	Update(user models.User) error
}

type userRepository struct {
	deps *deps.Dependencies
}

func NewUserRepository(deps *deps.Dependencies) UserRepository {
	return &userRepository{deps: deps}
}

func (ur userRepository) Create(req ent.UserCreateRequest, solarZ int64, brandType bool) error {
	ctx := context.Background()

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Company:      req.Company,
		PanelNumber:  req.PanelNumber,
		PotInstalled: req.PotInstalled,
		City:         req.City,
		Brand:        req.Brand,
		BrandType:    brandType,
		UserInverter: req.UserInverter,
		PassInverter: req.PassInverter,
	}
	if solarZ != 0 {
		user.PanelId = solarZ
	}

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		ur.deps.Logger.Debug().Err(err).Msgf("Error marshal values to insert DB")
		return fmt.Errorf("Erro encontrado, tente novamente mais tarde")
	}

	_, err = ur.deps.DB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      item,
	})
	if err != nil {
		ur.deps.Logger.Debug().Err(err).Msgf("Error inserting values DB")
		return fmt.Errorf("Erro encontrado, tente novamente mais tarde")
	}

	return nil
}

func (ur userRepository) GetByEmail(email string) (models.User, error) {
	key, err := attributevalue.MarshalMap(map[string]string{"Email": email})
	if err != nil {
		ur.deps.Logger.Warn().Err(err).Msgf("error creating key: %w", err)
		return models.User{}, fmt.Errorf("Erro encontrado, tente novamente mais tarde")
	}

	ctx := context.Background()
	out, err := ur.deps.DB.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key:       key,
	})
	if err != nil {
		ur.deps.Logger.Warn().Err(err).Msgf("erro no GetItem: %w", err)
		return models.User{}, fmt.Errorf("Erro encontrado, tente novamente mais tarde")
	}

	if out.Item == nil {
		return models.User{}, nil
	}

	var user models.User
	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		ur.deps.Logger.Warn().Err(err).Msgf("erro deserializando: %w", err)
		return models.User{}, fmt.Errorf("Erro encontrado, tente novamente mais tarde")
	}

	return user, nil
}

func (ur userRepository) Update(user models.User) error {
	ctx := context.Background()
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		ur.deps.Logger.Debug().Err(err).Msgf("Error marshal values to insert DB")
		return fmt.Errorf("Erro encontrado, tente novamente mais tarde")
	}

	_, err = ur.deps.DB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      item,
	})
	if err != nil {
		ur.deps.Logger.Debug().Err(err).Msgf("Error inserting values DB")
		return fmt.Errorf("Erro encontrado, tente novamente mais tarde")
	}

	return nil
}
