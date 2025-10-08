package middleware

import (
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	instance *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return getLabelValue(field)
	})

	return &Validator{
		instance: validate,
	}
}

func (v *Validator) Validate(r *http.Request, data any) (err error) {
	err = v.instance.Struct(data)
	if err != nil {
		for _, ve := range err.(validator.ValidationErrors) {
			return ve
		}
	}
	return nil
}

func getLabelValue(field reflect.StructField) string {
	tags := []string{"header", "json", "form", "path"}
	label := ""

	for _, tag := range tags {
		label = field.Tag.Get(tag)
		if label != "" {
			return label
		}
	}
	return ""
}
