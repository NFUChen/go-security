package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"go-security/internal"
	"html/template"
	"time"
)

type ResetPasswordClaims struct {
	UserID             uint    `json:"user_id"`
	ExpirationDuration float64 `json:"exp"`
}

func (claims *ResetPasswordClaims) Validate() error {
	if claims.ExpirationDuration < float64(time.Now().Unix()) {
		return internal.TokenExpired
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
	return &UserResetPasswordService{
		SmtpService: smtpService,
		UserService: userService,
		AuthService: authService,
		OtpService:  otpService,
	}
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
	return nil, internal.TokenInvalid
}

func (service *UserResetPasswordService) extractResetPasswordClaims(claims *jwt.MapClaims) (*ResetPasswordClaims, error) {
	var resetPasswordClaims ResetPasswordClaims
	purpose, ok := (*claims)["purpose"].(string)
	if !ok || purpose != string(PurposeResetPassword) {
		return nil, fmt.Errorf("invalid or missing 'purpose' claim")
	}
	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'user_id' claim")
	}
	expiration, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'exp' claim")
	}
	resetPasswordClaims.UserID = uint(userID)
	resetPasswordClaims.ExpirationDuration = expiration
	return &resetPasswordClaims, nil
}

func (service *UserResetPasswordService) doResetPassword(ctx context.Context, userId uint, newPassword string) error {
	user, err := service.UserService.FindUserByID(ctx, userId)
	if err != nil {
		return internal.UserNotFound
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
		return internal.ResetPasswordNotMatched
	}

	claims, err := service.parseResetPasswordClaims(token)
	if err != nil {
		return err
	}

	if err := service.OtpService.VerifyOtp(claims.UserID, PurposeResetPassword, otpCode); err != nil {
		return err
	}

	if err := claims.Validate(); err != nil {
		return err
	}
	return service.doResetPassword(ctx, claims.UserID, newPassword)
}

func (service *UserResetPasswordService) issueResetPasswordToken(email string) (string, error) {
	user, err := service.UserService.FindUserByEmail(context.Background(), email)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"purpose": string(PurposeResetPassword),
		"user_id": user.ID,
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
	}
	return service.AuthService.IssueJsonWebToken(&claims), nil
}

func (service *UserResetPasswordService) SendResetPasswordEmail(email string) (string, error) {
	token, err := service.issueResetPasswordToken(email)
	claims, err := service.parseResetPasswordClaims(token)
	if err != nil {
		return "", err
	}
	user, err := service.UserService.FindUserByID(context.Background(), claims.UserID)
	if err != nil {
		return "", err
	}
	otp := service.OtpService.GenerateOtp(claims.UserID, PurposeResetPassword)

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
	service.SmtpService.SendEmail(message)
	return token, nil

}
