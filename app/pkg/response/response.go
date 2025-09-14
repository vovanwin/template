package response

import (
	"app/config"
	"app/internal/shared/validator"
	"app/pkg/utils"
	"context"
	"errors"
	"github.com/go-chi/render"

	"github.com/vovanwin/platform/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

type Pagination struct {
	Limit      int64 `json:"limit,omitempty"`
	Page       int64 `json:"page,omitempty"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int64 `json:"total_pages"`
	IsFirst    bool  `json:"is_first"`
	IsLast     bool  `json:"is_last"`
}
type ResponseWithPagination struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Error struct {
	Data   interface{} `json:"data,omitempty"`
	Status string      `json:"status"`
}

type ResponseHandler struct {
	isProduction bool
}

func NewResponseHandler(config *config.Config) ResponseHandler {
	return ResponseHandler{
		isProduction: config.IsProduction(),
	}
}

// ErrorResponse для 400, 404 коды ошибки, и кастомное сообщение об ошибке.
func (h ResponseHandler) ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var errMsg string
	status := http.StatusBadRequest
	switch {
	case errors.Is(err, utils.ErrNotFound):
		status = http.StatusNotFound
		render.Status(r, status)
		render.JSON(w, r, Error{
			Status: "Error",
			Data:   err.Error(),
		})
		return
	case errors.Is(err, utils.ErrForbidden):
		logger.Warn(context.Background(), "нет доступа", zap.Error(err))
		status = http.StatusForbidden
		render.Status(r, status)
		render.JSON(w, r, Error{
			Status: "Error",
			Data:   err.Error(),
		})
		return
	case errors.Is(err, utils.ErrUnauthorized):
		logger.Warn(context.Background(), "Не авторизован", zap.Error(err))
		status = http.StatusUnauthorized
		render.Status(r, status)
		render.JSON(w, r, Error{
			Status: "Error",
			Data:   err.Error(),
		})
		return
	case errors.Is(err, utils.ErrValidation):
		logger.Warn(context.Background(), "Ошибка валидации", zap.Error(err))
		ValidationErrorResponse(w, r, err)
		return
	default:
		logger.Error(context.Background(), "ошибка запроса", zap.Error(err))
		status = http.StatusBadRequest
	}

	if h.isProduction {
		errMsg = "error"
	} else {
		errMsg = err.Error()
	}

	render.Status(r, status)
	render.JSON(w, r, Error{
		Status: "Error",
		Data:   errMsg,
	})
}

// ValidationErrorResponse Возвращает ошибки валидации, код ответа 422.
func ValidationErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// Создаем словарь для хранения ошибок валидации
	data := make(map[string]interface{})

	// Проверяем, является ли ошибка оберткой WrappedValidationError
	var validationErrs validator.WrappedValidationError
	if errors.As(err, &validationErrs) {
		// Если да, то извлекаем ошибки валидации
		for _, fieldErr := range validationErrs.ValidationErrors {
			if fieldErr.Error() == "" {
				continue
			}
			data[fieldErr.Field] = fieldErr.Error()
		}
	} else {
		// Если ошибка не связана с валидацией, добавляем общую ошибку
		data["error"] = err.Error()
	}

	render.Status(r, http.StatusUnprocessableEntity)
	render.JSON(w, r, Error{
		Status: "ValidationError",
		Data:   data,
	})
}

// SuccessResponse для кодов 200 и 204, успешные запросы.
func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
