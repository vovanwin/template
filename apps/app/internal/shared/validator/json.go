package validator

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

const _maxSize = 10240 // Установите максимальный размер JSON

// CustomValidator is a custom validator for JSON
// Custom validation function for checking the size of JSON
func (cv *CustomValidator) ValidateJSONSize(ctx context.Context, fl validator.FieldLevel) bool {
	// Преобразуем поле в JSON
	jsonValue, err := json.Marshal(fl.Field().Interface())
	if err != nil {
		return false
	}

	return len(jsonValue) <= _maxSize
}
