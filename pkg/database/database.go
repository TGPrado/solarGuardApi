package database

import (
	"context"
	"fmt"
	"log"

	secrets "github.com/TGPrado/GuardIA/config"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func connect(secrets *secrets.Config) (*dynamodb.Client, error) {
	if secrets.App.Environment != "dev" {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(secrets.DB.Region),
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar configuração do SDK: %w", err)
		}

		return dynamodb.NewFromConfig(cfg), nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(secrets.DB.Region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://db:8000"}, nil
			},
		)),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Credenciais fixas para ambiente local",
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("Erro ao carregar a configuração do SDK: %w", err)
	}
	return dynamodb.NewFromConfig(cfg), nil
}

func NewDynamoClient(secrets *secrets.Config, l *zerolog.Logger) (*dynamodb.Client, error) {
	ctx := context.Background()

	if secrets.DB.Endpoint == "" {
		return connect(secrets)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(secrets.DB.Region),
	)

	if err != nil {
		l.Debug().Err(err).Msgf("Erro conectando no banco de dados.")
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func CreateTableUsers(client *dynamodb.Client, l *zerolog.Logger) error {
	ctx := context.Background()

	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String("Users"),
	})

	if err == nil {
		log.Println("Tabela 'Users' já existe!")
		return nil
	}

	_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Email"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String("Users"),
		BillingMode: types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(25),
			WriteCapacityUnits: aws.Int64(25),
		},
	})
	if err != nil {
		l.Debug().Err(err).Msgf("Erro conectando no banco de dados.")
		return err
	}
	log.Println("Tabela 'Users' criada com sucesso!")
	return nil
}
