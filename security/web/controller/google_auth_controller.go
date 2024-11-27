package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/security/service"
	"go-security/security/service/oauth"
	"net/http"
	"time"
)

type GoogleAuthController struct {
	SecurityConfig    *service.SecurityConfig
	Router            *echo.Group
	GoogleAuthService *oauth.GoogleAuthService
}

func NewGoogleAuthController(routerGroup *echo.Group, googleAuthService *oauth.GoogleAuthService, securityConfig *service.SecurityConfig) *GoogleAuthController {
	return &GoogleAuthController{
		Router:            routerGroup,
		GoogleAuthService: googleAuthService,
		SecurityConfig:    securityConfig,
	}
}

func (controller *GoogleAuthController) RegisterRoutes() {
	controller.Router.POST("/public/google/login", controller.RegisterAndLogin)
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
	WriteCookie(&ctx, CookieName, token, expiration)
	return ctx.NoContent(http.StatusOK)
}
