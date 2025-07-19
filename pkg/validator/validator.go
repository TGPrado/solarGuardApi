package cfgValidator

import (
	"log"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ptBRTranslations "github.com/go-playground/validator/v10/translations/pt_BR"
)

func InitializeValidator(locale string) (*validator.Validate, ut.Translator) {
	validate := validator.New()

	enLocale := en.New()
	ptBRLocale := pt_BR.New()
	uni := ut.New(enLocale, enLocale, ptBRLocale)

	trans, ok := uni.GetTranslator(locale)
	if !ok {
		log.Fatalf("Language not supported: %s", locale)
	}

	switch locale {
	case "en":
		err := enTranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			log.Fatalf("Error registering english translation: %s", err)
		}
	case "pt_BR":
		err := ptBRTranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			log.Fatalf("Error registering portuguese translation: %s", err)
		}
	default:
		log.Fatalf("Language not supported: %s", locale)
	}

	return validate, trans
}

func TranslateErrors(err error, translator ut.Translator) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Translate(translator))
	}
	return errors
}
