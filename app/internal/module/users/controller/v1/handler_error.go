package usersv1

import (
	"context"
	api "github.com/vovanwin/template/internal/module/users/controller/gen"
	"net/http"
)

func (i Implementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{
		StatusCode: http.StatusBadRequest,
		Response: api.Error{
			Code:    400,
			Message: err.Error(),
		},
	}
}
