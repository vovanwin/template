package usersv1

import (
	api "app/internal/module/users/controller/gen"
	"context"
	"net/http"
)

func (i Implementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	if i.config.IsProduction() {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Code:    400,
				Message: "bad request",
			},
		}
	}

	return &api.ErrorStatusCode{
		StatusCode: http.StatusBadRequest,
		Response: api.Error{
			Code:    400,
			Message: err.Error(),
		},
	}
}
