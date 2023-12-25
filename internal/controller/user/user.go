package user

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"template/internal/repository/repoerrs"
	"template/internal/service/user"
	"template/pkg/utils"
)

type UserController struct {
	UserService user.UserService
	logger      *slog.Logger
}

func InitIndexController(e *echo.Echo, UserService user.UserService, logger *slog.Logger) {
	controller := &UserController{
		UserService: UserService,
		logger:      logger,
	}
	e.GET("/health", controller.first)
	e.GET("/", controller.first)
}

func (i *UserController) first(c echo.Context) error {
	userData, err := i.UserService.GetByID(c.Request().Context(), 1)

	if errors.Is(err, repoerrs.ErrNotFound) {
		utils.NewErrorResponse(c, http.StatusNotFound, "Не найдено")
		return err
	}

	return c.JSON(http.StatusOK, userData)
}
