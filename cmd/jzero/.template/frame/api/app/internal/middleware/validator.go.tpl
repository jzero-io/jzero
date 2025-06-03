package middleware

import (
	"net/http"
	"reflect"

	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type Validator struct {
	instance *validator.Validate
	trans    unTrans.Translator
}

func NewValidator() *Validator {
	validate := validator.New()
	uni := unTrans.New(zh_Hans_CN.New())

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return getLabelValue(field)
	})

	var err error
	// // register validation functions for custom validation
	// err = validate.RegisterValidation("customValidation", func(fl validator.FieldLevel) bool {
	//	return false
	// })

	trans, _ := uni.GetTranslator("zh_Hans_CN")
	err = zhTrans.RegisterDefaultTranslations(validate, trans)
	logx.Must(err)

	// // register custom validation error message
	// err = validate.RegisterTranslation("customValidation", trans, registerTranslator("customValidation", "自定义错误消息"), translate)
	// logx.Must(err)

	return &Validator{
		instance: validate,
		trans:    trans,
	}
}

func (v *Validator) Validate(r *http.Request, data any) (err error) {
	err = v.instance.Struct(data)
	if err != nil {
		for _, ve := range err.(validator.ValidationErrors) {
			if v.trans != nil {
				return errors.New(ve.Translate(v.trans))
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
			return label
		}
	}
	return ""
}

// func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
//	return func(trans unTrans.Translator) error {
//		if err := trans.Add(tag, msg, false); err != nil {
//			return err
//		}
//		return nil
//	}
// }

// func translate(trans unTrans.Translator, fe validator.FieldError) string {
//	msg, err := trans.T(fe.Tag(), fe.Field())
//	if err != nil {
//		panic(fe.(error).Error())
//	}
//	if len(strings.Split(fe.Namespace(), ".")) >= 2 {
//		return fmt.Sprintf("%s%s", strings.Split(fe.Namespace(), ".")[1], msg)
//	}
//	return msg
// }
