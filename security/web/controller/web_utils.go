package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/security/service"
	web "go-security/security/web/middleware"
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

func extractUserClaims(ctx echo.Context) (*service.UserClaims, error) {
	user := ctx.Get("user")
	claims, ok := user.(*service.UserClaims)
	if !ok {
		return nil, web.RoleKeyRequired
	}
	return claims, nil
}
