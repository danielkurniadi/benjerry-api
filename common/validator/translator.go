package validator

import (
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

var validate *validator.Validate
var trans ut.Translator

// NewValidator returns singleton validator.Validate
// which is threadsafe to be shared and used in concurrent
func NewValidator() (*validator.Validate, ut.Translator) {
	return validate, trans
}

func init() {
	validate = validator.New()
	localTrans := en.New()

	universal := ut.New(localTrans, localTrans)
	trans, _ = universal.GetTranslator("en")

	log.Println("app: Registering validate translator")

	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatal("error in registering validation translator:", err)
	}

	// Register validation field error to be translated/reformated
	registerTranslation(validate, trans, "required", "{0} is a required field")
	registerTranslation(validate, trans, "min", "{0} is less than minimmum required length")
	registerTranslation(validate, trans, "max", "{0} exceeded maximum character length")
	registerTranslation(validate, trans, "alpha", "{0} must contain only alphabet")
	registerTranslation(validate, trans, "numeric", "{0} must contain only numeric")
	registerTranslation(validate, trans, "alphanum", "{0} must contain only alphanumeric")
	registerTranslation(validate, trans, "ascii", "{0} must only contain ascii character")
	registerTranslation(validate, trans, "file", "{0} must be a valid unix file path")
	registerTranslation(validate, trans, "alpha", "{0} must only contain alphabet")
}

func registerTranslation(v *validator.Validate, trans ut.Translator, tag string, message string) {
	_ = v.RegisterTranslation(tag, trans,
		func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		},
	)
}
