package auth

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	authSer "template/internal/domain/user/service/auth"
	"template/pkg/utils/response"
	"template/pkg/validator"
)

type login struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

func (i *UserController) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input login

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		slog.Error("failed to decode request body", err)
		response.ErrorResponse(w, r, http.StatusBadRequest, "failed to decode request")
		return
	}

	if err := validator.NewCustomValidator().Validate(input); err != nil {
		slog.Error("invalid request", err)
		response.ValidationErrorResponse(w, r, err)
		return
	}

	token, err := i.AuthService.GetTokens(ctx, authSer.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		slog.Error("Ошибка запроса: ", err)
		response.ErrorResponse(w, r, http.StatusBadRequest, "Неверный логин или пароль")
		return
	}
	response.SuccessResponse(w, r, http.StatusOK, token)
}
