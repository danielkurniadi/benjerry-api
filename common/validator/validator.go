package validator

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// ValidationError ...
type ValidationError struct {
	err    error
	msg    string
	fields []string
}

// NewValidationError ...
func NewValidationError(err error) *ValidationError {
	if err == nil {
		return &ValidationError{err: nil, msg: "success"}
	}
	switch err.(type) {
	case *json.UnmarshalTypeError:
		verr, _ := err.(*json.UnmarshalTypeError)
		return fromJSONTypeError(verr)

	case validator.ValidationErrors:
		errs, _ := err.(validator.ValidationErrors)
		return fromValidatorError(errs)
	default:
		return &ValidationError{
			err: err,
			msg: err.Error(),
		}
	}
}

func fromJSONTypeError(err *json.UnmarshalTypeError) *ValidationError {
	msg := err.Field + " field must be of type " + err.Type.String()

	ve := &ValidationError{
		err:    err,
		msg:    msg,
		fields: []string{err.Field},
	}
	return ve
}

func fromValidatorError(errs validator.ValidationErrors) *ValidationError {
	var (
		msg       string
		fieldErrs []string
	)
	for _, e := range errs {
		fieldErrs = append(fieldErrs, e.Translate(trans))
	}
	msg = strings.Join(fieldErrs, "; ")

	ve := &ValidationError{
		err:    errs,
		msg:    msg,
		fields: fieldErrs,
	}
	return ve
}

func (ve *ValidationError) Error() string { return ve.err.Error() }

// Message provides client-friendly message
func (ve *ValidationError) Message() string {
	return ve.msg
}

var (
	validate *validator.Validate
	trans    ut.Translator
	isInit   bool
)

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

	setupRegisteredTranslations(validate, trans)
	isInit = true
}

// ValidateStruct ,,,
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	verr := NewValidationError(err)
	return verr
}

// ValidateVar ..
func ValidateVar(field interface{}, tag string) error {
	err := validate.Var(field, tag)
	if err == nil {
		return nil
	}

	verr := NewValidationError(err)
	return verr
}

// ValidateJSON ...
func ValidateJSON(field io.Reader, s interface{}) error {
	unmarshalErr := json.NewDecoder(field).Decode(s)
	if unmarshalErr != nil {
		return NewValidationError(unmarshalErr)
	}
	return nil
}

// DecodeAndValidateJSON unmarshal json string into struct
// and if the struct fields are tagged with `validate`, will
// also perform struct validation
func DecodeAndValidateJSON(field io.Reader, s interface{}) error {
	unmarshalErr := json.NewDecoder(field).Decode(s)
	if unmarshalErr != nil {
		return NewValidationError(unmarshalErr)
	}
	return ValidateStruct(s)
}
