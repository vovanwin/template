package validator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// PasswordValidationErrors представляет ошибки валидации пароля.
type PasswordValidationErrors []string

// Error возвращает текстовое представление ошибок валидации пароля.
func (pve PasswordValidationErrors) Error() string {
	return strings.Join(pve, ", ")
}

const (
	passwordMinLength = 16
	passwordMaxLength = 32
	passwordMinLower  = 1
	passwordMinUpper  = 1
	passwordMinDigit  = 1
	passwordMinSymbol = 1
)

var (
	lowerCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[a-z]{%d,}`, passwordMinLower))
	upperCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[A-Z]{%d,}`, passwordMinUpper))
	digitRegexp     = regexp.MustCompile(fmt.Sprintf(`[0-9]{%d,}`, passwordMinDigit))
	symbolRegexp    = regexp.MustCompile(fmt.Sprintf(`[!@#$%%^&*]{%d,}`, passwordMinSymbol))
)

func (cv *CustomValidator) passwordValidate(ctx context.Context, fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		cv.passwdErr = fmt.Errorf("поле %s должно быть строкой", fl.FieldName())
		return false
	}

	fieldValue := fl.Field().String()
	var errors PasswordValidationErrors

	// Проверка длины пароля
	if len(fieldValue) < passwordMinLength || len(fieldValue) > passwordMaxLength {
		errors = append(errors, fmt.Sprintf("в поле %s должно быть между %d и %d символами", fl.FieldName(), passwordMinLength, passwordMaxLength))
	}

	// Массив для проверки регулярных выражений
	var rules = []struct {
		regex *regexp.Regexp
		err   string
	}{
		{lowerCaseRegexp, fmt.Sprintf("поле %s должно содержать как минимум %d строчную букву", fl.FieldName(), passwordMinLower)},
		{upperCaseRegexp, fmt.Sprintf("поле %s должно содержать как минимум %d прописную букву", fl.FieldName(), passwordMinUpper)},
		{digitRegexp, fmt.Sprintf("поле %s должно содержать как минимум %d цифру", fl.FieldName(), passwordMinDigit)},
		{symbolRegexp, fmt.Sprintf("поле %s должно содержать как минимум %d специальный символ", fl.FieldName(), passwordMinSymbol)},
	}

	// Применение проверок и накопление ошибок
	for _, rule := range rules {
		if !rule.regex.MatchString(fieldValue) {
			errors = append(errors, rule.err)
		}
	}

	if len(errors) > 0 {
		cv.passwdErr = errors
		return false
	}

	return true
}
