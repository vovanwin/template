package response

import (
	"github.com/go-chi/render"

	"net/http"
)

type Pagination struct {
	Limit      int64 `json:"limit,omitempty;query:limit"`
	Page       int64 `json:"page,omitempty;query:page"`
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

//
//// ResponsePagination Структура для ответа с пагинацией,limit и page берутся из контекста
//// который записывается в middleware
//func ResponsePagination(w http.ResponseWriter, r *http.Request, statusCode int, data map[string]any) {
//	c := r.Context()
//
//	limit, _ := c.Value(framework.Limit).(int64)
//	page, _ := c.Value(framework.Page).(int64)
//	totalRows := data["count"].(int64)
//	totalPages := int64(math.Ceil(float64(totalRows) / float64(limit)))
//	dataRes := data["data"].(interface{})
//
//	res := ResponseWithPagination{
//		Data: dataRes,
//		Pagination: Pagination{
//			Limit:      limit,
//			Page:       page,
//			TotalRows:  totalRows,
//			TotalPages: totalPages,
//			IsFirst:    page <= 1,
//			IsLast:     totalPages == page,
//		},
//	}
//
//	render.Status(r, statusCode)
//	render.JSON(w, r, res)
//}

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
func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
