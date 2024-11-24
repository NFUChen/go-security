package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-security/security"
	. "go-security/security/repository"
	"html/template"
	"time"
)

type UserVerificationClaims struct {
	UserID             uint    `json:"user_id"`
	ExpirationDuration float64 `json:"exp"`
}

type UserVerificationService struct {
	SmtpService ISmtpService
	UserService *UserService
	AuthService *AuthService
	OtpService  *OtpService
}

func (service *UserVerificationService) issueVerificationToken(user *User) (string, error) {
	if user.IsVerified {
		return "", security.UserAlreadyVerified
	}

	claims := jwt.MapClaims{
		"purpose": string(PurposeGuestEmailVerification),
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}

	_jwt := service.AuthService.IssueJsonWebToken(&claims)

	return _jwt, nil
}

func (service *UserVerificationService) SendVerificationEmail(ctx context.Context, email string) (string, error) {
	user, err := service.UserService.FindUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	_template, err := template.New("email_verification").Parse(EMAIL_VERIFICATION_HTML_TEMPLATE)
	if err != nil {
		return "", err
	}
	otp := service.OtpService.GenerateOtp(user.ID, PurposeGuestEmailVerification)

	var buffer bytes.Buffer
	emailTemplate := NewEmailTemplate(user.Name, otp.Code, service.SmtpService.GetSmtpConfig().CompanyName)
	if err := _template.Execute(&buffer, emailTemplate); err != nil {
		return "", err
	}

	emailContent := buffer.String()

	message := service.SmtpService.CreateNewMessage(email, "Email Verification", emailContent, ContentTypeHtml)
	service.SmtpService.SendEmail(message)
	token, err := service.issueVerificationToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
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
	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		return nil, security.TokenInvalid
	}
	verificationClaims.UserID = uint(userID)
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

	if err := service.OtpService.VerifyOtp(claims.UserID, PurposeGuestEmailVerification, otpCode); err != nil {
		return err
	}
	user, err := service.UserService.FindUserByID(context.Background(), claims.UserID)

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
