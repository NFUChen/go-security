package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func writeCookie(c *echo.Context, key string, value string, duration time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Expires = time.Now().Add(duration)
	cookie.Path = "/"
	(*c).SetCookie(cookie)
}
