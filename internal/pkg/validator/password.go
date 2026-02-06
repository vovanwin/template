package validator

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	passwordMinLength = 8
	passwordMaxLength = 32
	passwordMinLower  = 1
	passwordMinUpper  = 1
	passwordMinDigit  = 1
	passwordMinSymbol = 1
)

var (
	lengthRegexp    = regexp.MustCompile(fmt.Sprintf(`^.{%d,%d}$`, passwordMinLength, passwordMaxLength))
	lowerCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[a-z]{%d,}`, passwordMinLower))
	upperCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[A-Z]{%d,}`, passwordMinUpper))
	digitRegexp     = regexp.MustCompile(fmt.Sprintf(`[0-9]{%d,}`, passwordMinDigit))
	symbolRegexp    = regexp.MustCompile(fmt.Sprintf(`[!@#$%%^&*]{%d,}`, passwordMinSymbol))
)

func (cv *CustomValidator) passwordValidate(fl validator.FieldLevel) bool {
	// проверить что поле это строка
	if fl.Field().Kind() != reflect.String {
		cv.passwdErr = fmt.Errorf("поле %s должно быть строкой", fl.FieldName())
		return false
	}

	// get the value of the field
	fieldValue := fl.Field().String()

	// check regexp matching
	if ok := lengthRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("в поле %s должно быть между %d и %d символами", fl.FieldName(), passwordMinLength, passwordMaxLength)
		return false
	}
	//else if ok = lowerCaseRegexp.MatchString(fieldValue); !ok {
	//	cv.passwdErr = fmt.Errorf("field %s must contain at least %d lowercase letter(s)", fl.FieldName(), passwordMinLower)
	//	return false
	//} else if ok = upperCaseRegexp.MatchString(fieldValue); !ok {
	//	cv.passwdErr = fmt.Errorf("field %s must contain at least %d uppercase letter(s)", fl.FieldName(), passwordMinUpper)
	//	return false
	//} else if ok = digitRegexp.MatchString(fieldValue); !ok {
	//	cv.passwdErr = fmt.Errorf("field %s must contain at least %d digit(s)", fl.FieldName(), passwordMinDigit)
	//	return false
	//} else if ok = symbolRegexp.MatchString(fieldValue); !ok {
	//	cv.passwdErr = fmt.Errorf("field %s must contain at least %d special character(s)", fl.FieldName(), passwordMinSymbol)
	//	return false
	//}

	return true
}
