package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type Error struct {
	Data   interface{} `json:"data,omitempty"`
	Status string      `json:"status"`
}

// ValidationErrorResponse Возвращаем когда ошибка валидации, код ответа 422
func ValidationErrorResponse(w http.ResponseWriter, r *http.Request, errs error) {
	render.Status(r, http.StatusUnprocessableEntity)
	render.JSON(w, r, Error{
		Status: "ValidationError",
		Data:   errs.Error(),
	})
}

// ErrorResponse для 400, 404 коды ошибки, и кстомное сообщение об ошибке
func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, msg string) {
	render.Status(r, status)
	render.JSON(w, r, Error{
		Status: "Error",
		Data:   msg,
	})
}

// SuccessResponse для кодов 200 и 204, успешные запросы
func OkResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}

// SuccessResponse для кодов 200 и 204, успешные запросы
func NoContentResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, http.StatusNoContent)
	render.JSON(w, r, data)
}
