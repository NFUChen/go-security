package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go-security/security/service"
	web "go-security/security/web/middleware"
)

type EmailRateLimitedController struct {
	*AuthController
	UserService *service.UserService
	Router      *echo.Group
}

func NewEmailRateLimitedController(limitedRRouter *echo.Group, userService *service.UserService, authController *AuthController) *EmailRateLimitedController {
	return &EmailRateLimitedController{AuthController: authController, UserService: userService, Router: limitedRRouter}
}

func (controller *EmailRateLimitedController) RegisterRoutes() {
	log.Info().Msgf("Registering email rate limited routes...")
	supperAdmin, err := controller.UserService.GetRoleByName(context.TODO(), service.RoleSuperAdmin)
	if err != nil {
		log.Fatal().Msgf("Failed to find role: %v", err)
	}
	controller.Router.POST("/private/send-verification-email-by-user-id", web.RoleRequired(supperAdmin, controller.AdminSendVerificationEmailByUserID))

	controller.Router.POST("/public/send-reset-password-email", controller.SendResetPasswordEmail)
	controller.Router.POST("/private/send-verification-email-by-token", controller.SendVerificationEmailByToken)
}
