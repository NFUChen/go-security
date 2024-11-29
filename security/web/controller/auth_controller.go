package controller

import (
	"github.com/labstack/echo/v4"
	"go-security/security/service"
	"net/http"
	"time"
)

type AuthController struct {
	Router                   *echo.Group
	SecurityConfig           *service.SecurityConfig
	AuthService              *service.AuthService
	UserResetPasswordService *service.UserResetPasswordService
	UserVerificationService  *service.UserVerificationService
	UserService              *service.UserService
}

func NewAuthController(routerGroup *echo.Group, authService *service.AuthService, userResetPasswordService *service.UserResetPasswordService, userVerificationService *service.UserVerificationService, userService *service.UserService, securityConfig *service.SecurityConfig) *AuthController {
	return &AuthController{
		Router:                   routerGroup,
		AuthService:              authService,
		UserResetPasswordService: userResetPasswordService,
		UserVerificationService:  userVerificationService,
		UserService:              userService,
		SecurityConfig:           securityConfig,
	}
}

const (
	CookieName = "jwt"
)

func (controller *AuthController) RegisterRoutes() {

	controller.Router.POST("/public/register", controller.RegisterUser)
	controller.Router.POST("/public/login", controller.Login)
	controller.Router.GET("/private/current-user", controller.GetUser)
	controller.Router.POST("/public/issue-reset-password-token", controller.IssueResetPasswordToken)

	controller.Router.POST("/public/reset-password", controller.ResetPassword)

	controller.Router.GET("/private/issue-verification-token", controller.IssueVerificationToken)
	controller.Router.GET("/private/is-admin-pushed-email-verification", controller.IsAdminPushedEmailVerification)
	controller.Router.POST("/private/verify-email", controller.VerifyEmail)

	controller.Router.GET("/private/logout", controller.Logout)
	controller.Router.GET("/private/redirect-url", controller.GetRedirectURL)
}

func (controller *AuthController) GetRedirectURL(ctx echo.Context) error {
	user, _ := ExtractUserClaims(ctx)
	isAdmin := controller.AuthService.IsUserAdmin(user.RoleName)
	var redirectUrl string
	if isAdmin {
		redirectUrl = controller.SecurityConfig.AdminRedirectUrl
	} else {
		redirectUrl = controller.SecurityConfig.ClientRedirectUrl
	}

	return ctx.JSON(http.StatusOK, map[string]string{"redirectURL": redirectUrl})
}

func (controller *AuthController) GetUser(ctx echo.Context) error {
	userClaims, err := ExtractUserClaims(ctx)
	if err != nil {
		WriteCookie(&ctx, CookieName, "", -1*time.Hour)
		return err
	}

	return ctx.JSON(http.StatusOK, userClaims)
}

func (controller *AuthController) IsAdminPushedEmailVerification(ctx echo.Context) error {
	userClaims, err := ExtractUserClaims(ctx)
	if err != nil {
		return err
	}
	isAskingVerification := controller.UserVerificationService.IsAdminAskingForVerification(userClaims.ID)

	return ctx.JSON(http.StatusOK, map[string]bool{"is_admin_pushed": isAskingVerification})

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

	registeredUser, err := controller.AuthService.RegisterUserAsGuest(ctx.Request().Context(), user.UserName, user.Email, user.Password, service.PlatformSelf, nil)
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
	WriteCookie(&ctx, CookieName, token, 24*time.Hour)
	return ctx.NoContent(http.StatusOK)
}

func (controller *AuthController) Logout(ctx echo.Context) error {
	WriteCookie(&ctx, CookieName, "", -1*time.Hour)
	return ctx.NoContent(http.StatusOK)
}

// can be also used for resending.
func (controller *AuthController) SendResetPasswordEmail(ctx echo.Context) error {
	var resetPasswordSchema struct {
		Token string `json:"token"`
	}
	if err := ctx.Bind(&resetPasswordSchema); err != nil {
		return err
	}
	token, err := controller.UserResetPasswordService.SendResetPasswordEmail(ctx.Request().Context(), resetPasswordSchema.Token)
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
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Password has been reset"})
}

func (controller *AuthController) IssueVerificationToken(ctx echo.Context) error {
	// for verification, we can assume that the user is already logged in, but the email is not verified
	userClaims, err := ExtractUserClaims(ctx)
	if err != nil {
		return err
	}
	token, err := controller.UserVerificationService.IssueVerificationToken(ctx.Request().Context(), userClaims.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"token": token})
}

func (controller *AuthController) IssueResetPasswordToken(ctx echo.Context) error {
	var schema struct {
		Email string `json:"email"`
	}
	if err := ctx.Bind(&schema); err != nil {
		return err
	}
	token, err := controller.UserResetPasswordService.IssueResetPasswordToken(ctx.Request().Context(), schema.Email)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"token": token})
}

func (controller *AuthController) SendVerificationEmailByToken(ctx echo.Context) error {
	var schema struct {
		Token string `json:"token"`
	}
	if err := ctx.Bind(&schema); err != nil {
		return err
	}
	err := controller.UserVerificationService.SendVerificationEmailByToken(ctx.Request().Context(), schema.Token)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusAccepted)
}

func (controller *AuthController) AdminSendVerificationEmailByUserID(ctx echo.Context) error {
	var schema struct {
		UserID uint `json:"user_id"`
	}
	if err := ctx.Bind(&schema); err != nil {
		return err
	}

	err := controller.UserVerificationService.SendVerificationEmailByUserID(ctx.Request().Context(), schema.UserID, true)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusAccepted)
}

func (controller *AuthController) VerifyEmail(ctx echo.Context) error {
	var verificationSchema struct {
		Token string `json:"token"`
		Otp   string `json:"otp"`
	}
	if err := ctx.Bind(&verificationSchema); err != nil {
		return err
	}
	err := controller.UserVerificationService.VerifyEmail(verificationSchema.Token, verificationSchema.Otp)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Email has been verified"})
}
