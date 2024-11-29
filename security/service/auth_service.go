package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"go-security/security"
	. "go-security/security/repository"
	"golang.org/x/crypto/bcrypt"
	"slices"
	"time"
)

const (
	CookieName = "jwt"
)

type SecurityConfig struct {
	Secret                string   `yaml:"secret" json:"-"`
	ExcludedRoutePrefixes []string `yaml:"excluded_routes_prefixes"`
	AdminRedirectUrl      string   `yaml:"admin_redirect_url"`
	ClientRedirectUrl     string   `yaml:"client_redirect_url"`
}

type UserClaims struct {
	ID                 uint    `json:"id"`
	UserName           string  `json:"user_name"`
	RoleName           string  `json:"role"`
	RoleIndex          uint    `json:"role_index"`
	ExpirationDuration float64 `json:"exp"`
	IsVerified         bool    `json:"is_verified"`
}

func NewUserClaims(userID uint, userName string, roleName string, roleIndex uint, expiration float64, isVerified bool) *UserClaims {
	return &UserClaims{
		ID:                 userID,
		UserName:           userName,
		RoleName:           roleName,
		RoleIndex:          roleIndex,
		ExpirationDuration: expiration,
		IsVerified:         isVerified,
	}
}

func (claims *UserClaims) Validate() error {
	if claims.ExpirationDuration < float64(time.Now().Unix()) {
		return security.TokenExpired
	}
	return nil
}

type AuthService struct {
	Secret      string
	UserService *UserService
}

func NewAuthService(userService *UserService, secret string) *AuthService {
	authService := &AuthService{
		Secret:      secret,
		UserService: userService,
	}

	return authService
}

func (service *AuthService) PostConstruct() {

}

func (service *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := service.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", security.UserNotFound
	}
	err = service.VerifyPassword(password, user.Password)
	if err != nil {
		return "", security.UserPasswordNotMatched
	}
	return service.IssueLoginToken(user, time.Hour)
}

func (service *AuthService) IssueJsonWebToken(claims *jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(service.Secret))
	log.Info().Msgf("Issue Token: %v", tokenString)
	return tokenString
}

func (service *AuthService) IsUserAdmin(userRoleName string) bool {
	admins := []string{RoleSuperAdmin, RoleAdmin}
	if slices.Contains(admins, userRoleName) {
		return true
	}
	return false
}

func (service *AuthService) IssueLoginToken(user *User, expiration time.Duration) (string, error) {

	claims := jwt.MapClaims{
		"user_name":   user.Name,
		"id":          user.ID,
		"role_name":   user.Role.Name,
		"role_index":  user.Role.RoleIndex,
		"exp":         time.Now().Add(expiration).Unix(),
		"is_verified": user.IsVerified,
	}
	return service.IssueJsonWebToken(&claims), nil
}

func (service *AuthService) ExtractUserClaims(claims *jwt.MapClaims) (*UserClaims, error) {

	userID, ok := (*claims)["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'id' claim")
	}
	userName, ok := (*claims)["user_name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'user_name' claim")
	}

	roleName, ok := (*claims)["role_name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'role' claim")
	}

	roleIndex, ok := (*claims)["role_index"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'role_index' claim")
	}

	expiration, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'exp' claim")
	}

	isVerified, ok := (*claims)["is_verified"].(bool)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'is_verified' claim")
	}

	userClaims := NewUserClaims(uint(userID), userName, roleName, uint(roleIndex), expiration, isVerified)

	return userClaims, nil
}

func (service *AuthService) DecodeJsonWebTokenWithSecret(rawToken string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (service *AuthService) DecodeJsonWebToken(rawToken string) (*jwt.Token, error) {
	return service.DecodeJsonWebTokenWithSecret(rawToken, []byte(service.Secret))
}

func (service *AuthService) ParseUserClaims(tokenString string) (*UserClaims, error) {
	_jwt, err := service.DecodeJsonWebToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := _jwt.Claims.(jwt.MapClaims); ok && _jwt.Valid {
		userClaims, err := service.ExtractUserClaims(&claims)
		if err != nil {
			return nil, err
		}

		return userClaims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func (service *AuthService) GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (service *AuthService) VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (service *AuthService) NewUser(name string, email string, password string, role *UserRole, platform *Platform, externalID *string) (*User, error) {

	user := &User{
		Name:       name,
		Email:      email,
		Password:   password,
		RoleID:     role.ID,
		IsVerified: false,
		Platform:   *platform,
		ExternalID: externalID,
	}

	err := user.Validate()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *AuthService) RegisterUserAsGuest(ctx context.Context, name string, email string, password string, platformName PlatformType, externalID *string) (*User, error) {
	return service.RegisterUser(ctx, name, email, password, platformName, externalID, RoleGuest)
}

func (service *AuthService) RegisterUser(ctx context.Context, name string, email string, password string, platformType PlatformType, externalID *string, userRole string) (*User, error) {
	existingUser, err := service.UserService.GetUserByEmail(ctx, email)
	if err == nil {
		log.Info().Msgf("User already exists: %v", existingUser)
		return existingUser, security.UserAlreadyExists
	}
	hashedPassword, err := service.GenerateHashPassword(password)
	if err != nil {
		return nil, err
	}
	role, err := service.UserService.GetRoleByName(ctx, userRole)
	if err != nil {
		return nil, err
	}

	platform, err := service.UserService.GetPlatformByName(ctx, platformType)
	if err != nil {
		return nil, err
	}

	user, err := service.NewUser(name, email, hashedPassword, role, platform, externalID)
	if err != nil {
		return nil, err
	}

	return service.UserService.SaveUser(ctx, user)
}
