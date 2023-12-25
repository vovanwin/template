package httpserver

import (
	"github.com/labstack/echo/v4"
)

func RequestIdTest(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		c.Set("RequestId", c.Response().Header().Get(echo.HeaderXRequestID))

		return next(c)
	}
}
