package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-security/security"
	. "go-security/security/service"
	"net/http"
	"net/url"
	"time"
)

type LineConfig struct {
	ClientID           string `yaml:"client_id"`
	ClientSecret       string `yaml:"client_secret"`
	ChannelSecret      string `yaml:"channel_secret"`
	ChannelAccessToken string `yaml:"channel_access_token"`
	RedirectURI        string `yaml:"redirect_uri"`
}

type LineAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

type LineUser struct {
	Issuer                string   `json:"iss"`     // The entity that issued the token
	Subject               string   `json:"sub"`     // The unique identifier of the user
	Audience              string   `json:"aud"`     // The intended audience of the token (usually the application or service)
	ExpirationTime        int      `json:"exp"`     // The expiration time of the token in Unix timestamp format
	IssuedAtTime          int      `json:"iat"`     // The time the token was issued in Unix timestamp format
	AuthenticationMethods []string `json:"amr"`     // The list of methods used to authenticate the user (e.g., password, multi-factor authentication)
	UserName              string   `json:"name"`    // The full name of the user
	ProfilePictureURL     string   `json:"picture"` // The URL of the user's profile picture
	Email                 string   `json:"email"`   // The email address of the user
}

type LineLoginService struct {
	LineConfig  *LineConfig
	AuthService *AuthService
	UserService *UserService
	Http        *http.Client
}

func NewLineLoginService(authService *AuthService, userService *UserService, lineConfig *LineConfig) *LineLoginService {
	return &LineLoginService{
		AuthService: authService,
		UserService: userService,
		LineConfig:  lineConfig,
		Http:        &http.Client{},
	}
}

func (service *LineLoginService) GetLineLoginUrl() string {
	baseUrl := "https://access.line.me/oauth2/v2.1/authorize?"
	state := uuid.New().String()                     // Unique state to prevent CSRF
	scope := url.QueryEscape("openid profile email") // Permissions required
	return fmt.Sprintf(
		"%sresponse_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		baseUrl, service.LineConfig.ClientID, service.LineConfig.RedirectURI, state, scope,
	)
}

func (service *LineLoginService) GetLineAuthResponse(code string) (*LineAuthResponse, error) {
	tokenUrl := "https://api.line.me/oauth2/v2.1/token"
	// Create the request body as x-www-form-urlencoded format
	param := url.Values{}
	param.Set("grant_type", "authorization_code")
	param.Set("redirect_uri", service.LineConfig.RedirectURI)
	param.Set("client_id", service.LineConfig.ClientID)
	param.Set("client_secret", service.LineConfig.ClientSecret)
	param.Set("code", code)
	// Create the HTTP request
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBufferString(param.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// Set headers for x-www-form-urlencoded content type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Send the POST request
	resp, err := service.Http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into LineAuthResponse struct
	var lineAuthResponse LineAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&lineAuthResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &lineAuthResponse, nil
}

func (service *LineLoginService) VerifyIDToken(idToken string) (*LineUser, error) {
	// Define the URL for verification
	verifyURL := "https://api.line.me/oauth2/v2.1/verify"
	// Prepare the request body
	data := url.Values{}
	data.Set("id_token", idToken)
	data.Set("client_id", service.LineConfig.ClientID)

	// Create the HTTP request
	req, err := http.NewRequest("POST", verifyURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var lineUser LineUser
	if err := json.NewDecoder(resp.Body).Decode(&lineUser); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &lineUser, nil
}

func (service *LineLoginService) RegisterAndLogin(ctx context.Context, user *LineUser) (string, error) {
	targetUser, err := service.UserService.GetUserByEmail(ctx, user.Email)
	if targetUser == nil && errors.Is(err, security.UserNotFound) {
		targetUser, err = service.AuthService.RegisterUserAsGuest(ctx, user.UserName, user.Email, user.Subject, PlatformLine, &user.Subject)
	}
	if err != nil {
		return "", err
	}
	if err := service.UserService.ActivateUser(ctx, targetUser); err != nil {
		return "", err
	}

	expirationTime := time.Until(time.Unix(int64(user.ExpirationTime), 0))
	return service.AuthService.IssueLoginToken(targetUser, expirationTime)
}
