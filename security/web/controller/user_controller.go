package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"go-security/security/repository"
	"go-security/security/service"
	web "go-security/security/web/middleware"
	"net/http"
)

type UserController struct {
	Router               *echo.Group
	UserService          *service.UserService
	ResetPasswordService *service.UserResetPasswordService
	VerificationService  *service.UserVerificationService
}

func NewUserController(routerGroup *echo.Group, userService *service.UserService, resetPasswordService *service.UserResetPasswordService, verificationService *service.UserVerificationService) *UserController {
	return &UserController{
		Router:               routerGroup,
		UserService:          userService,
		ResetPasswordService: resetPasswordService,
		VerificationService:  verificationService,
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
