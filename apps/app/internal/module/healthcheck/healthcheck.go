package healthcheck

import (
	"app/pkg/buildinfo"
	"app/pkg/response"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Implementation struct {
}

func Controller(r *chi.Mux) {
	controller := &Implementation{}

	r.Group(func(r chi.Router) {
		r.Route("/api/v1/healthcheck", func(r chi.Router) {
			r.Get("/", controller.healthcheck)
		})
		r.Route("/version", func(r chi.Router) {
			r.Get("/", controller.version)
		})
	})
}

func (i *Implementation) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (i *Implementation) version(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(w, r, http.StatusOK, buildinfo.BuildInfo)
}
