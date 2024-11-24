package web

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go-security/security/service"
	"net/http"
	"strings"
)

var (
	LoginRequired    = errors.New("LoginRequired")
	PermissionDenied = errors.New("PermissionDenied")
	RoleKeyRequired  = errors.New("RoleKeyRequired")
)

type AuthMiddleware struct {
	AuthService    *service.AuthService
	ExcludedRoutes []string
}

func NewAuthMiddleware(authService *service.AuthService, excludedRoutes []string) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService:    authService,
		ExcludedRoutes: excludedRoutes,
	}
}

func (middleware *AuthMiddleware) AuthMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		urlPath := ctx.Request().URL.Path
		for _, excludedRoute := range middleware.ExcludedRoutes {
			if strings.HasPrefix(urlPath, excludedRoute) {
				ctx.Logger().Infof("Bypass auth middleware excluded route: %s", urlPath)
				return next(ctx)
			}
		}

		cookie, err := ctx.Cookie(service.CookieName)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, LoginRequired)
		}
		userClaims, err := middleware.AuthService.ParseUserClaims(cookie.Value)
		if err != nil {
			return err
		}

		ctx.Set("user", userClaims)
		return next(ctx)
	}
}

func RoleRequired(requiredRoleIndex uint, next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user")
		castedUser, ok := user.(*service.UserClaims)
		if !ok {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": RoleKeyRequired.Error(),
			})
		}

		if castedUser.RoleIndex < requiredRoleIndex {
			log.Warn().Msgf("User %s with role %s has no permission to access this resource", castedUser.UserName, castedUser.RoleName)
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": PermissionDenied.Error(),
			})
		}
		return next(ctx)
	}
}
