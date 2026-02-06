package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	v         *validator.Validate
	passwdErr error
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.RegisterValidation("password", cv.passwordValidate)
	if err != nil {
		panic(err)
	}

	err = v.RegisterValidation("IsISO8601Date", cv.iso8601Validate)
	if err != nil {
		panic(err)
	}

	return cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.v.Struct(i)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]

		return cv.newValidationError(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
	}
	return nil
}

func (cv *CustomValidator) newValidationError(field string, value interface{}, tag string, param string) error {
	switch tag {
	case "required":
		return fmt.Errorf("поле %s является обязательным", field)
	case "email":
		return fmt.Errorf("поле %s должно быть валидным Емейл адресом", field)
	case "password":
		return cv.passwdErr
	case "min":
		return fmt.Errorf("поле %s должно быть не меньше чем %s символов", field, param)
	case "max":
		return fmt.Errorf("поле %s должно быть не больше чем %s символов", field, param)
	case "IsISO8601Date":
		return cv.passwdErr
	default:
		return fmt.Errorf("поле %s невалидно", field)
	}
}
