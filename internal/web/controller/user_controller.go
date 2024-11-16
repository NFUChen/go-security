package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"go-security/internal/repository"
	"go-security/internal/service"
	web "go-security/internal/web/middleware"
	"net/http"
)

type UserController struct {
	Router      *echo.Group
	UserService *service.UserService
}

func NewUserController(routerGroup *echo.Group, userService *service.UserService) *UserController {
	return &UserController{
		Router:      routerGroup,
		UserService: userService,
	}
}

func (controller *UserController) RegisterRoutes() {
	adminRole, err := controller.UserService.FindRoleByName(context.Background(), repository.RoleAdmin)
	if err != nil {
		panic(err)
	}
	controller.Router.GET("/private/user", web.RoleRequired(adminRole.RoleIndex, controller.GetUser))
}

func (controller *UserController) GetUser(ctx echo.Context) error {
	users, err := controller.UserService.FindAllUsers(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, users)
}
