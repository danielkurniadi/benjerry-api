package validator

import (
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// List of translation for validation error message
var (
	// message for checking required field
	RequiredValidateMessage = "{0} is a required field"

	// message for checking minimum/maximum length
	MinValidateMessage = "{0} is less than minimmum required length"
	MaxValidateMessage = "{0} is less than minimmum required length"

	// message for checking string characters
	AlphabetValidateMessage = "{0} must contain only alphabetical"
	AlphaNumValidateMessage = "{0} must contain only alphanumeric"
	ASCIIValidateMessage    = "{0} must only contain ascii character"
	NumericValidateMessage  = "{0} must contain only numeric"

	// message for checking string is a file path
	FileValidateMessage = "{0} must be a valid unix file path"
)

// setupRegisteredTranslations registers validation field
// error message to be translated/reformated to the above listed message
// such that it is more client-friendly
func setupRegisteredTranslations(validate *validator.Validate, trans ut.Translator) {
	registerTranslation(validate, trans, "required", RequiredValidateMessage)
	registerTranslation(validate, trans, "min", MinValidateMessage)
	registerTranslation(validate, trans, "max", MaxValidateMessage)

	registerTranslation(validate, trans, "alpha", AlphabetValidateMessage)
	registerTranslation(validate, trans, "numeric", NumericValidateMessage)
	registerTranslation(validate, trans, "alphanum", AlphaNumValidateMessage)
	registerTranslation(validate, trans, "ascii", ASCIIValidateMessage)

	registerTranslation(validate, trans, "file", FileValidateMessage)
}

// registerTranslation is a helper to register translated field error
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
