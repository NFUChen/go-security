package oauth

import (
	"context"
	"errors"
	"fmt"
	"go-security/security"
	. "go-security/security/service"
	"time"
)

// GoogleUser represents the user context information returned from a Google authentication flow.
type GoogleUser struct {
	SubjectIdentifier    string  `json:"sub"`               // The subject identifier, a unique identifier for the user.
	Email                string  `json:"email"`             // The email address of the user.
	IsEmailVerified      bool    `json:"email_verified"`    // Whether the email address has been verified.
	FirstName            string  `json:"given_name"`        // The given name (first name) of the user.
	LastName             string  `json:"family_name"`       // The family name (last name) of the user.
	ProfilePictureSource *string `json:"picture,omitempty"` // The URL of the user's profile picture, optional.
	IssuedAt             int64   `json:"iat"`               // The time the ID token was issued (Unix epoch seconds).
	Expiration           int64   `json:"exp"`               // The time the ID token expires (Unix epoch seconds).
	Issuer               string  `json:"iss"`               // The issuer identifier, typically the URL of Google's OAuth 2.0 Authorization Server.
	Audience             string  `json:"aud"`               // The audience (recipient of the ID token).
	NonceValue           *string `json:"nonce,omitempty"`   // A string to associate a client session with the ID token, optional.
}

func (ctx *GoogleUser) FullName() string {
	return fmt.Sprintf("%s %s", ctx.FirstName, ctx.LastName)
}

type GoogleAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type GoogleAuthService struct {
	AuthConfig  *GoogleAuthConfig
	AuthService *AuthService
	UserService *UserService
}

func NewGoogleAuthConfig(clientID string, clientSecret string) *GoogleAuthConfig {
	return &GoogleAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func NewGoogleAuthService(authConfig *GoogleAuthConfig, authService *AuthService, userService *UserService) *GoogleAuthService {
	return &GoogleAuthService{
		AuthConfig:  authConfig,
		AuthService: authService,
		UserService: userService,
	}
}

func (service *GoogleAuthService) RegisterAndLogin(ctx context.Context, user *GoogleUser) (string, error) {
	targetUser, err := service.UserService.GetUserByEmail(ctx, user.Email)
	if targetUser == nil && errors.Is(err, security.UserNotFound) {
		targetUser, err = service.AuthService.RegisterUserAsGuest(ctx, user.FullName(), user.Email, user.SubjectIdentifier, PlatformGoogle, &user.SubjectIdentifier)
		if err != nil {
			return "", err
		}
	}
	if targetUser == nil {
		return "", security.UserNotFound
	}

	expirationTime := time.Until(time.Unix(user.Expiration, 0))
	if expirationTime <= 0 {
		return "", security.TokenExpired
	}
	if user.IsEmailVerified {
		if err := service.AuthService.UserService.ActivateUser(ctx, targetUser); err != nil {
			return "", err
		}
	}
	// Issue a login token for the user.
	token, err := service.AuthService.IssueLoginToken(targetUser, expirationTime)
	if err != nil {
		return "", err
	}

	return token, nil
}
