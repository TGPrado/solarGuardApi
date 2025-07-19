package helpers

import (
	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	cfgValidator "github.com/TGPrado/GuardIA/pkg/validator"
	"github.com/gin-gonic/gin"
)

func ValidateInput(req interface{}, c *gin.Context, deps *deps.Dependencies) (string, bool) {
	if err := c.ShouldBindJSON(&req); err != nil {
		return err.Error(), false
	}

	if err := deps.Validator.Struct(req); err != nil {
		errors := cfgValidator.TranslateErrors(err, deps.Translator)
		return errors[0], false
	}

	return "", true
}
