package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"go-security/security"
	"html/template"
	"time"
)

type ResetPasswordClaims struct {
	ID                 uint    `json:"id"`
	ExpirationDuration float64 `json:"exp"`
}

func (claims *ResetPasswordClaims) Validate() error {
	if claims.ExpirationDuration < float64(time.Now().Unix()) {
		return security.TokenExpired
	}
	return nil
}

type UserResetPasswordService struct {
	SmtpService ISmtpService
	UserService *UserService
	AuthService *AuthService
	OtpService  *OtpService
}

func NewUserResetPasswordService(smtpService ISmtpService, userService *UserService, authService *AuthService, otpService *OtpService) *UserResetPasswordService {
	service := &UserResetPasswordService{
		SmtpService: smtpService,
		UserService: userService,
		AuthService: authService,
		OtpService:  otpService,
	}
	fmt.Printf("")
	return service

}

func (service *UserResetPasswordService) parseResetPasswordClaims(token string) (*ResetPasswordClaims, error) {
	_jwt, err := service.AuthService.DecodeJsonWebToken(token)
	if err != nil {
		return nil, err
	}

	if claims, ok := _jwt.Claims.(jwt.MapClaims); ok && _jwt.Valid {
		resetPasswordClaims, err := service.extractResetPasswordClaims(&claims)
		if err != nil {
			return nil, err
		}
		return resetPasswordClaims, nil
	}
	return nil, security.TokenInvalid
}

func (service *UserResetPasswordService) extractResetPasswordClaims(claims *jwt.MapClaims) (*ResetPasswordClaims, error) {
	var resetPasswordClaims ResetPasswordClaims
	purpose, ok := (*claims)["purpose"].(string)
	if !ok || purpose != string(PurposeResetPassword) {
		return nil, fmt.Errorf("invalid or missing 'purpose' claim, getting %s, expects %v", purpose, PurposeResetPassword)
	}
	userID, ok := (*claims)["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'id' claim")
	}
	expiration, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'exp' claim")
	}
	resetPasswordClaims.ID = uint(userID)
	resetPasswordClaims.ExpirationDuration = expiration
	return &resetPasswordClaims, nil
}

func (service *UserResetPasswordService) doResetPassword(ctx context.Context, userId uint, newPassword string) error {
	user, err := service.UserService.GetUserByID(ctx, userId)
	if err != nil {
		return security.UserNotFound
	}
	hashedPassword, err := service.AuthService.GenerateHashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := service.UserService.ResetUserPassword(ctx, user, hashedPassword); err != nil {
		return err
	}
	return nil
}

func (service *UserResetPasswordService) ResetPassword(ctx context.Context, token string, otpCode string, newPassword string, confirmedPassword string) error {
	if newPassword != confirmedPassword {
		return security.ResetPasswordNotMatched
	}
	log.Info().Msgf("Reset Password: %v", token)
	claims, err := service.parseResetPasswordClaims(token)
	if err != nil {
		log.Warn().Msgf("Failed to parse reset password claims: %v", err)
		return err
	}

	log.Info().Msgf("Claims: %v", claims)
	log.Info().Msgf("Begin to verify OTP: %v", otpCode)
	if err := service.OtpService.VerifyOtp(claims.ID, PurposeResetPassword, otpCode); err != nil {
		return err
	}

	if err := claims.Validate(); err != nil {
		return err
	}
	return service.doResetPassword(ctx, claims.ID, newPassword)
}

func (service *UserResetPasswordService) IssueResetPasswordToken(ctx context.Context, email string) (string, error) {
	user, err := service.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user.Platform.Name != string(PlatformSelf) {
		return "", security.SelfPlatformRequiredForPasswordReset
	}

	claims := jwt.MapClaims{
		"purpose": string(PurposeResetPassword),
		"id":      user.ID,
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
	}
	return service.AuthService.IssueJsonWebToken(&claims), nil
}

func (service *UserResetPasswordService) SendResetPasswordEmail(context context.Context, token string) (string, error) {

	claims, err := service.parseResetPasswordClaims(token)
	if err != nil {
		return "", err
	}

	user, err := service.UserService.GetUserByID(context, claims.ID)
	if err != nil {
		return "", err
	}

	otp := service.OtpService.GenerateOtp(claims.ID, PurposeResetPassword)

	subject := "Reset Password"
	_template, err := template.New("reset_password_email").Parse(RESET_PASSWORD_EMAIL_HTML_TEMPLATE)
	if err != nil {
		return "", err
	}
	emailTemplate := NewEmailTemplate(user.Name, otp.Code, service.SmtpService.GetSmtpConfig().CompanyName)
	var buffer bytes.Buffer
	if err := _template.Execute(&buffer, emailTemplate); err != nil {
		log.Info().Msgf("Failed to execute template: %v", err)
	}
	body := buffer.String()
	message := service.SmtpService.CreateNewMessage(user.Email, subject, body, ContentTypeHtml)
	if err := service.SmtpService.SendEmail(message); err != nil {
		return "", err
	}
	return token, nil

}
