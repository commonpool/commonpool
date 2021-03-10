package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

var DefaultValidator = NewValidator()
var enLocale = en.New()
var UniversalTranslator = ut.New(enLocale, enLocale)
var DefaultTranslator, _ = UniversalTranslator.GetTranslator("en")

type CustomValidator interface {
	Validate() error
}

func NewValidator() *Validator {

	var (
		validate = validator.New()
	)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = validate.RegisterValidation("notblank", validators.NotBlank, true)

	_ = enTranslations.RegisterDefaultTranslations(validate, DefaultTranslator)

	_ = validate.RegisterTranslation("required", DefaultTranslator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = validate.RegisterTranslation("notblank", DefaultTranslator, func(ut ut.Translator) error {
		return ut.Add("notblank", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notblank", fe.Field())
		return t
	})

	return &Validator{
		validator: validate,
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) RegisterStruct(fn func(sl validator.StructLevel), types ...interface{}) {
	v.validator.RegisterStructValidation(fn, types...)
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	if c1, ok := i.(CustomValidator); ok {
		return c1.Validate()
	}
	return nil
}
