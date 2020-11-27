package router

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var Trans ut.Translator

var DefaultValidator = NewValidator()

func NewValidator() *Validator {

	var (
		validate = validator.New()
		enLocale = en.New()
		uni      = ut.New(enLocale, enLocale)
		trans, _ = uni.GetTranslator("en")
	)

	Trans = trans

	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	return &Validator{
		validator: validate,
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
