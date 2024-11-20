package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/internal/service"
	"net/http"
	"time"
)

type AuthController struct {
	AuthService              *service.AuthService
	UserResetPasswordService *service.UserResetPasswordService
	UserVerificationService  *service.UserVerificationService
	UserService              *service.UserService
	Router                   *echo.Group
}

const (
	CookieName = "jwt"
)

func (controller *AuthController) RegisterRoutes() {
	controller.Router.POST("/public/register", controller.RegisterUser)
	controller.Router.POST("/public/login", controller.Login)

	controller.Router.POST("/public/send-reset-password-email", controller.SendResetPasswordEmail)
	controller.Router.POST("/public/reset-password", controller.ResetPassword)

	controller.Router.POST("/public/send-verification-email", controller.SendVerificationEmail)
	controller.Router.POST("/public/verify-email", controller.VerifyEmail)

	controller.Router.GET("/private/logout", controller.Logout)

}

func (controller *AuthController) RegisterUser(ctx echo.Context) error {
	var user struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	registeredUser, err := controller.AuthService.RegisterUser(ctx.Request().Context(), user.UserName, user.Email, user.Password)
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

func (controller *AuthController) SendResetPasswordEmail(ctx echo.Context) error {
	var resetPasswordSchema struct {
		Email string `json:"email"`
	}
	if err := ctx.Bind(&resetPasswordSchema); err != nil {
		return err
	}
	token, err := controller.UserResetPasswordService.SendResetPasswordEmail(resetPasswordSchema.Email)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"token": token})
}

func (controller *AuthController) ResetPassword(ctx echo.Context) error {
	var resetPasswordSchema struct {
		Token             string `json:"token"`
		OtpCode           string `json:"otp_code"`
		NewPassword       string `json:"new_password"`
		ConfirmedPassword string `json:"confirmed_password"`
	}
	if err := ctx.Bind(&resetPasswordSchema); err != nil {
		return err
	}
	err := controller.UserResetPasswordService.ResetPassword(
		ctx.Request().Context(),
		resetPasswordSchema.Token,
		resetPasswordSchema.OtpCode,
		resetPasswordSchema.NewPassword,
		resetPasswordSchema.ConfirmedPassword,
	)

	if err != nil {
		return err
	}
	return ctx.String(http.StatusOK, "Password reset successfully")
}

func (controller *AuthController) SendVerificationEmail(ctx echo.Context) error {
	var emailSchema struct {
		Email string `json:"email"`
	}
	if err := ctx.Bind(&emailSchema); err != nil {
		return err
	}
	token, err := controller.UserVerificationService.SendVerificationEmail(ctx.Request().Context(), emailSchema.Email)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"token": token})
}

func (controller *AuthController) VerifyEmail(ctx echo.Context) error {
	var verificationSchema struct {
		Token   string `json:"token"`
		OtpCode string `json:"otp_code"`
	}
	if err := ctx.Bind(&verificationSchema); err != nil {
		return err
	}
	err := controller.UserVerificationService.VerifyEmail(verificationSchema.Token, verificationSchema.OtpCode)
	if err != nil {
		return err
	}
	return ctx.String(http.StatusOK, "Email verified successfully")
}

func NewAuthController(routerGroup *echo.Group, authService *service.AuthService, userService *service.UserService) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
		Router:      routerGroup,
	}
}
