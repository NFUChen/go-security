package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"runtime/debug"
)

func ErrorMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			stacks := string(debug.Stack())
			c.Logger().Error(stacks)
			fmt.Println(stacks)

			// Return a custom error response
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}
		return err
	}
}
