package deps

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/rs/zerolog"
)

type Dependencies struct {
	Handler    *gin.Engine
	Logger     *zerolog.Logger
	DB         *dynamodb.Client
	Validator  *validator.Validate
	Translator ut.Translator
}
