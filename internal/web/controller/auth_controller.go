package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/internal/service"
	"net/http"
	"time"
)

type AuthController struct {
	AuthService *service.AuthService
	UserService *service.UserService
	Router      *echo.Group
}

const (
	CookieName = "jwt"
)

func (controller *AuthController) RegisterRoutes() {
	controller.Router.POST("/public/register", controller.RegisterUser)
	controller.Router.POST("/public/login", controller.Login)
	controller.Router.GET("/private/logout", controller.Logout)

}

func (controller *AuthController) RegisterUser(ctx echo.Context) error {
	var user struct {
		Name     string `json:"user_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	registeredUser, err := controller.AuthService.RegisterUser(ctx.Request().Context(), user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusAccepted, registeredUser)
}

func (controller *AuthController) Login(ctx echo.Context) error {
	var loginCredential struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Bind(&loginCredential); err != nil {
		return err
	}

	token, err := controller.AuthService.Login(ctx.Request().Context(),
		loginCredential.Email, loginCredential.Password)
	if err != nil {
		return err
	}
	writeCookie(&ctx, CookieName, token, 24*time.Hour)
	return ctx.String(http.StatusOK, "Login successfully")
}

func (controller *AuthController) Logout(ctx echo.Context) error {
	writeCookie(&ctx, CookieName, "", -1*time.Hour)
	return ctx.String(http.StatusOK, "Logout successfully")
}

func NewAuthController(routerGroup *echo.Group, authService *service.AuthService, userService *service.UserService) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
		Router:      routerGroup,
	}
}
