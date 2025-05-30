package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(r *http.Request, data any) (err error) {
	validate := validator.New()
	uni := unTrans.New(zh_Hans_CN.New())

	// register validation functions for custom validation
	//err = validate.RegisterValidation("customValidation", func(fl validator.FieldLevel) bool {
	//	return false
	//})

	trans, _ := uni.GetTranslator("zh_Hans_CN")
	err = zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return err
	}

	// register custom validation error message
	//err = validate.RegisterTranslation("customValidation", trans, registerTranslator("customValidation", "自定义错误消息"), translate)
	//if err != nil {
	//	return err
	//}

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return getLabelValue(field)
	})

	err = validate.Struct(data)
	if err != nil {
		for _, ve := range err.(validator.ValidationErrors) {
			if trans != nil {
				return errors.Errorf(ve.Translate(trans))
			}
			return ve
		}
	}
	return nil
}

func getLabelValue(field reflect.StructField) string {
	tags := []string{"label", "json", "form", "path"}
	label := ""

	for _, tag := range tags {
		label = field.Tag.Get(tag)
		if label != "" {
			break
		}
	}
	return ""
}

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans unTrans.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

func translate(trans unTrans.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	if len(strings.Split(fe.Namespace(), ".")) >= 2 {
		return fmt.Sprintf("%s%s", strings.Split(fe.Namespace(), ".")[1], msg)
	}
	return msg
}