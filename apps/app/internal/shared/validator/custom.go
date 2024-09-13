package validator

import (
	"app/internal/shared/validator/dbsqlc"
	"app/pkg/storage/postgres"
	"app/pkg/utils"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError представляет ошибку валидации.
type ValidationError struct {
	Field string      // Название поля с ошибкой.
	Value interface{} // Значение поля.
	Tag   string      // Тег валидации, вызвавший ошибку.
	Param string      // Дополнительные параметры валидации.
}

// Error возвращает текстовое представление ошибки валидации.
func (ve ValidationError) Error() string {
	switch ve.Tag {
	case "required":
		return fmt.Sprintf("Поле '%s' является обязательным.", ve.Field)
	case "email":
		return fmt.Sprintf("Поле '%s' должно содержать валидный адрес электронной почты.", ve.Field)
	case "password":
		return ve.Param
	case "eqfield":
		return fmt.Sprintf("Поле '%s' не соответствует полю '%s'.", ve.Field, ve.Param)
	case "min":
		return fmt.Sprintf("Поле '%s' должно быть не короче %s символов.", ve.Field, ve.Param)
	case "max":
		return fmt.Sprintf("Поле '%s' должно быть не длиннее %s символов.", ve.Field, ve.Param)
	case "oneof":
		return fmt.Sprintf("Поле '%s' должно быть одним из следующих значений: [%s].", ve.Field, ve.Param)
	case "IsISO8601Date":
		return fmt.Sprintf("Поле '%s' должно содержать валидную дату в формате ISO 8601.", ve.Field)
	case "isAllowDevicesValidate":
		return fmt.Sprintf("Устройства не существуют или нет прав на взаимодействие с ними  [%s]", ve.Param)
	case "exists":
		return fmt.Sprintf("Не найдено совпадений по идентификатору [%s]", ve.Param)
	case "unique":
		return fmt.Sprintf("Запись с таким значением уже существует [%s]", ve.Value)
	case "json_size":
		return fmt.Sprintf("Размер json не может быть больше [%d] символов", _maxSize)
	default:
		return fmt.Sprintf("поле %s невалидно", ve.Field)
	}
}

// ValidationErrors представляет массив ошибок валидации.
type ValidationErrors []ValidationError

// Error возвращает текстовое представление массива ошибок валидации.
func (ve ValidationErrors) Error() string {
	var errors []string
	for _, err := range ve {
		errors = append(errors, err.Error())
	}
	return strings.Join(errors, ", ")
}

// CustomValidator представляет кастомный валидатор.
type CustomValidator struct {
	v         *validator.Validate // Экземпляр валидатора.
	passwdErr error               // Ошибка пароля.
	pgx       *postgres.Postgres  // Подключение к базе данных PostgreSQL.
	sqlc      *dbsqlc.Queries     // Подключение к базе данных PostgreSQL.
}

// New создает новый экземпляр CustomValidator и регистрирует кастомные правила валидации.
func New(pgx *postgres.Postgres) *CustomValidator {
	v := validator.New()
	sqlc := dbsqlc.New(pgx.Pool)
	cv := &CustomValidator{v: v, pgx: pgx, sqlc: sqlc}

	// Регистрация функции для получения имени тега JSON.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Регистрация кастомных правил валидации.
	if err := v.RegisterValidationCtx("unique", cv.uniqueValidate); err != nil {
		panic(err)
	}
	if err := v.RegisterValidationCtx("password", cv.passwordValidate); err != nil {
		panic(err)
	}
	if err := v.RegisterValidationCtx("IsISO8601Date", cv.iso8601Validate); err != nil {
		panic(err)
	}

	if err := v.RegisterValidationCtx("exists", cv.existsValidate); err != nil {
		panic(err)
	}
	if err := v.RegisterValidationCtx("isAllowDevicesValidate", cv.isAllowDevicesValidate); err != nil {
		panic(err)
	}
	if err := v.RegisterValidationCtx("json_size", cv.ValidateJSONSize); err != nil {
		panic(err)
	}
	return cv
}

// WrappedValidationError оборачивает ошибки валидации с помощью глобальной переменной ErrValidation.
type WrappedValidationError struct {
	ValidationErrors ValidationErrors
}

// Error возвращает текстовое представление ошибок.
func (w WrappedValidationError) Error() string {
	return w.ValidationErrors.Error()
}

// Is позволяет использовать errors.Is для проверки типа ошибки.
func (w WrappedValidationError) Is(target error) bool {
	return errors.Is(utils.ErrValidation, target)
}

func (cv *CustomValidator) Validate(ctx context.Context, i interface{}) error {
	if i == nil {
		return fmt.Errorf("передан nil объект для валидации")
	}

	var allErrors ValidationErrors

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Slice {
		// Валидация для срезов
		for j := 0; j < v.Len(); j++ {
			elem := v.Index(j).Interface()
			cv.collectValidationErrors(ctx, elem, &allErrors)
		}
	} else {
		// Валидация для одиночного объекта
		cv.collectValidationErrors(ctx, i, &allErrors)
	}

	if len(allErrors) > 0 {
		return WrappedValidationError{ValidationErrors: allErrors}
	}

	return nil
}

// collectValidationErrors собирает ошибки валидации для переданного объекта
func (cv *CustomValidator) collectValidationErrors(ctx context.Context, obj interface{}, allErrors *ValidationErrors) {
	if err := cv.v.StructCtx(ctx, obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				*allErrors = append(*allErrors, ValidationError{
					Field: fieldErr.Field(),
					Value: fieldErr.Value(),
					Tag:   fieldErr.Tag(),
					Param: fieldErr.Param(),
				})
			}
		} else {
			*allErrors = append(*allErrors, ValidationError{
				Field: "unknown",
				Value: nil,
				Tag:   "unknown",
				Param: err.Error(),
			})
		}
	}

	// Проверка ошибок пароля
	if cv.passwdErr != nil {
		if passwordErrors, ok := cv.passwdErr.(PasswordValidationErrors); ok {
			for _, err := range passwordErrors {
				*allErrors = append(*allErrors, ValidationError{
					Field: "password",
					Value: nil,
					Tag:   "password",
					Param: err,
				})
			}
		} else {
			*allErrors = append(*allErrors, ValidationError{
				Field: "password",
				Value: nil,
				Tag:   "password",
				Param: cv.passwdErr.Error(),
			})
		}
		cv.passwdErr = nil
	}
}
