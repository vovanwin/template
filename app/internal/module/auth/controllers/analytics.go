package controller

import (
	"github.com/vovanwin/template/pkg/response"
	"net/http"
)

func (i *Implementation) TEST(w http.ResponseWriter, r *http.Request) {
	//i.service.GetLogin(r.Context())

	type name struct {
		Name string
	}
	response.OkResponse(w, r, name{
		Name: "asdas",
	})
}
