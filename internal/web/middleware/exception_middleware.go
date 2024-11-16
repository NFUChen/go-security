package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func ErrorMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			// Handle the error, log it or send a custom response
			c.Logger().Error(err)
			// Return a custom error response
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		return err
	}
}
