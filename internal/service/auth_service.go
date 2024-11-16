package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"go-security/internal"
	. "go-security/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	CookieName = "jwt"
)

type UserClaims struct {
	UserName           string  `json:"user_name"`
	RoleName           string  `json:"role"`
	RoleIndex          uint    `json:"role_index"`
	ExpirationDuration float64 `json:"exp"`
}

type AuthService struct {
	Secret      string
	UserService *UserService
	AllRoles    []*UserRole
}

func NewAuthService(userService *UserService, secret string) *AuthService {
	authService := &AuthService{
		Secret:      secret,
		UserService: userService,
	}
	authService.addBuiltinRoles()
	return authService
}

func (service *AuthService) addBuiltinRoles() {
	for _, role := range BuiltinRoles {
		_ = service.UserService.AddRole(context.Background(), &role)
	}
}

func (service *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := service.UserService.FindUserByEmail(ctx, email)
	if err != nil {
		return "", internal.UserNotFound
	}
	role, err := service.UserService.FindRoleByName(ctx, user.Role.Name)
	err = service.VerifyPassword(password, user.Password)
	if err != nil {
		return "", internal.UserPasswordNotMatched
	}
	return service.encodeJsonWebToken(user.Name, role, time.Hour)
}

func (service *AuthService) encodeJsonWebToken(userName string, role *UserRole, expiration time.Duration) (string, error) {

	userMapClaims := jwt.MapClaims{
		"user_name":  userName,
		"role_name":  role.Name,
		"role_index": role.RoleIndex,
		"exp":        time.Now().Add(expiration).Unix(),
	}
	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userMapClaims)

	// Sign the token with the secret key
	return token.SignedString([]byte(service.Secret))
}

func (service *AuthService) ExtractUserClaims(claims *jwt.MapClaims) (*UserClaims, error) {
	userClaims := UserClaims{}
	userName, ok := (*claims)["user_name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'user_name' claim")
	}
	userClaims.UserName = userName

	roleName, ok := (*claims)["role_name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'role' claim")
	}
	userClaims.RoleName = roleName

	roleIndex, ok := (*claims)["role_index"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'role_index' claim")
	}
	userClaims.RoleIndex = uint(roleIndex)

	expiration, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'exp' claim")
	}
	userClaims.ExpirationDuration = expiration

	return &userClaims, nil
}

func (service *AuthService) DecodeJsonWebToken(tokenString string) (*UserClaims, error) {
	// Parse the token with a function to validate the signature
	secretKey := []byte(service.Secret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
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

func (service *AuthService) NewUser(name, email, password, role string) (*User, error) {
	userRole, err := service.UserService.FindRoleByName(context.Background(), role)
	if err != nil {
		return nil, err
	}
	user := &User{
		Name:     name,
		Email:    email,
		Password: password,
		RoleID:   userRole.ID,
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *AuthService) RegisterUser(ctx context.Context, name string, email string, password string) (*User, error) {
	existingUser, err := service.UserService.FindUserByEmail(ctx, email)
	if err == nil {
		log.Info().Msgf("User already exists: %v", existingUser)
		return nil, internal.UserAlreadyExists
	}
	hashedPassword, err := service.GenerateHashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := service.NewUser(name, email, hashedPassword, RoleGuest)
	if err != nil {
		return nil, err
	}
	return service.UserService.SaveUser(ctx, user)
}
