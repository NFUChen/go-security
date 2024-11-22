package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/security/service/oauth"
	"net/http"
	"time"
)

type GoogleAuthController struct {
	Router            *echo.Group
	GoogleAuthService *oauth.GoogleAuthService
}

func (controller *GoogleAuthController) RegisterRoutes() {
	controller.Router.POST("/public/google/login", controller.RegisterAndLogin)
}

func NewGoogleAuthController(router *echo.Group, authService *oauth.GoogleAuthService) *GoogleAuthController {
	return &GoogleAuthController{Router: router, GoogleAuthService: authService}
}

func (controller *GoogleAuthController) RegisterAndLogin(ctx echo.Context) error {
	var googleUser oauth.GoogleUser
	if err := ctx.Bind(&googleUser); err != nil {
		return err
	}
	token, err := controller.GoogleAuthService.RegisterAndLogin(ctx.Request().Context(), &googleUser)
	if err != nil {
		return err
	}
	expiration := time.Until(time.Unix(googleUser.Expiration, 0))
	writeCookie(&ctx, CookieName, token, expiration)
	return ctx.String(http.StatusOK, "Google login successfully")
}
