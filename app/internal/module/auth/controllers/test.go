package controller

import (
	"github.com/vovanwin/template/pkg/response"
	"net/http"
)

func (i *Implementation) TEST(w http.ResponseWriter, r *http.Request) {
	login, err := i.service.GetLogin(r.Context())
	if err != nil {
		return
	}

	response.OkResponse(w, r, login)
}
