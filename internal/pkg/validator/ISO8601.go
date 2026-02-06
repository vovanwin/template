package validator

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	ISO8601DateRegex = regexp.MustCompile("^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$")
)

// Правило валидации для формата времени  iso8601 2023-04-17T01:20:00Z
func (cv *CustomValidator) iso8601Validate(fl validator.FieldLevel) bool {
	// Проверяем, если поле равно Null то возвращаем true в проверке
	isNil := fl.Field().IsZero()
	if isNil == true {
		return true
	}

	// проверить что поле это строка
	if fl.Field().Kind() != reflect.String {
		cv.passwdErr = fmt.Errorf("поле %s должно быть строкой", fl.FieldName())
		return false
	}

	// получить значение поля
	fieldValue := fl.Field().String()

	// check regexp matching
	if ok := ISO8601DateRegex.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("поле %s несоответствует формату времени ISO-8601  ", fl.FieldName())
		return false
	}

	return true
}
