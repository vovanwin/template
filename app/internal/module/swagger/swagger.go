package swagger

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/pkg/response"
	"github.com/vovanwin/template/pkg/validator"
	"log/slog"
	"net/http"
)

type Implementation struct{ config *config.Config }

type Param struct {
	Link string `json:"param" validate:"omitempty,url"`
}

func Controller(r *chi.Mux, config *config.Config) {
	controller := &Implementation{config: config}
	if config.Env != "prod" {
		r.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth("введите пароль и логин", map[string]string{"user": "111"}))
			r.Get("/api/v1/swagger", controller.swaggerNew)
			r.Get("/api/v1/swagger/file", controller.swaggerFile)
			r.Get("/api/v1/swagger/file/{file}", controller.swaggerFile)

		})
	}
}

// Новый сваггер
func (i *Implementation) swaggerNew(w http.ResponseWriter, r *http.Request) {
	var list map[string]string
	req := Param{
		Link: r.URL.Query().Get("link"),
	}
	if err := validator.NewCustomValidator().Validate(req); err != nil {
		slog.Error("invalid request", "err", err)
		response.ValidationErrorResponse(w, r, err)
		return
	}

	var baseUrl string
	if i.config.Env == "local" {
		baseUrl = "http://localhost:8080/api/v1/swagger/newFile"
	}
	if i.config.Env == "dev" {
		baseUrl = "http://localhost:8080/api/v1/swagger/newFile"
	}
	if req.Link == "" {
		req.Link = baseUrl
	}

	list = map[string]string{
		"Пользователи v1": "/api/v1/swagger/new?link=" + baseUrl + "/usersv1",
	}

	index := newIndexPage()
	err := index.handler(w, list, req)
	if err != nil {
		return
	}

}

func (i *Implementation) swaggerFile(w http.ResponseWriter, r *http.Request) {
	//file := chi.URLParam(r, "file")
	//var sw *openapi3.T
	//var err error
	//
	//switch file {
	//case "usersv1":
	//	sw, err = usersGenv1.GetSwagger()
	//}
	//
	//if err != nil {
	//	response.ErrorResponse(w, r, http.StatusBadRequest, "")
	//}
	//response.OkResponse(w, r, sw)

}
