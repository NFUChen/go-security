package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-security/security"
	"html/template"
	"time"
)

type UserVerificationClaims struct {
	ID                 uint    `json:"id"`
	ExpirationDuration float64 `json:"exp"`
}

type UserVerificationService struct {
	SmtpService ISmtpService
	UserService *UserService
	AuthService *AuthService
	OtpService  *OtpService
}

func (service *UserVerificationService) IssueVerificationToken(ctx context.Context, userID uint) (string, error) {
	user, err := service.UserService.FindUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.IsVerified {
		return "", security.UserAlreadyVerified
	}

	claims := jwt.MapClaims{
		"purpose": string(PurposeGuestEmailVerification),
		"id":      user.ID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}

	_jwt := service.AuthService.IssueJsonWebToken(&claims)

	return _jwt, nil
}

func (service *UserVerificationService) SendVerificationEmail(ctx context.Context, token string) error {
	claims, err := service.parseVerificationClaims(token)
	if err != nil {
		return err
	}

	user, err := service.UserService.FindUserByID(ctx, claims.ID)
	if err != nil {
		return err
	}

	_template, err := template.New("email_verification").Parse(EMAIL_VERIFICATION_HTML_TEMPLATE)
	if err != nil {
		return err
	}
	otp := service.OtpService.GenerateOtp(user.ID, PurposeGuestEmailVerification)

	var buffer bytes.Buffer
	emailTemplate := NewEmailTemplate(user.Name, otp.Code, service.SmtpService.GetSmtpConfig().CompanyName)
	if err := _template.Execute(&buffer, emailTemplate); err != nil {
		return err
	}

	emailContent := buffer.String()

	message := service.SmtpService.CreateNewMessage(user.Email, "Email Verification", emailContent, ContentTypeHtml)
	service.SmtpService.SendEmail(message)

	return nil
}

func (service *UserVerificationService) parseVerificationClaims(token string) (*UserVerificationClaims, error) {
	_jwt, err := service.AuthService.DecodeJsonWebToken(token)
	if err != nil {
		return nil, err
	}
	if claims, ok := _jwt.Claims.(jwt.MapClaims); ok && _jwt.Valid {
		verificationClaims, err := service.extractVerificationClaims(&claims)
		if err != nil {
			return nil, err
		}
		return verificationClaims, nil
	}
	return nil, security.TokenInvalid
}

func (service *UserVerificationService) extractVerificationClaims(claims *jwt.MapClaims) (*UserVerificationClaims, error) {
	var verificationClaims UserVerificationClaims
	purpose, ok := (*claims)["purpose"].(string)
	if !ok || purpose != string(PurposeGuestEmailVerification) {
		return nil, fmt.Errorf("invalid or missing 'purpose' claim, getting %s, expects %v", purpose, PurposeGuestEmailVerification)
	}
	userID, ok := (*claims)["id"].(float64)
	if !ok {
		return nil, security.TokenInvalid
	}
	verificationClaims.ID = uint(userID)
	exp, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, security.TokenInvalid
	}
	verificationClaims.ExpirationDuration = exp
	return &verificationClaims, nil
}

func (service *UserVerificationService) VerifyEmail(token string, otpCode string) error {
	claims, err := service.parseVerificationClaims(token)
	if err != nil {
		return err
	}

	if err := service.OtpService.VerifyOtp(claims.ID, PurposeGuestEmailVerification, otpCode); err != nil {
		return err
	}
	user, err := service.UserService.FindUserByID(context.Background(), claims.ID)

	return service.UserService.ActivateUser(context.Background(), user)
}

func NewUserVerificationService(smtpService ISmtpService, userService *UserService, authService *AuthService, otpService *OtpService) *UserVerificationService {
	return &UserVerificationService{
		SmtpService: smtpService,
		UserService: userService,
		AuthService: authService,
		OtpService:  otpService,
	}
}
