package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-security/security/service"
	"go-security/security/web"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

var emailRateLimitConfig = middleware.RateLimiterConfig{
	Skipper: middleware.DefaultSkipper,
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(4),    // 4 emails
			Burst:     4,                // Burst size equal to rate limit
			ExpiresIn: 10 * time.Minute, // Time window of 10 minutes
		},
	),
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		claims, ok := ctx.Get("user").(*service.UserClaims)
		if !ok {
			return ctx.RealIP(), nil
		}
		return strconv.Itoa(int(claims.ID)), nil
	},
	ErrorHandler: func(context echo.Context, err error) error {
		return context.JSON(http.StatusForbidden, map[string]string{"message": web.UnableToIdentifyUser.Error()})
	},
	DenyHandler: func(context echo.Context, identifier string, err error) error {
		return context.JSON(http.StatusTooManyRequests, map[string]string{"message": web.EmailRateLimitExceeded.Error()})
	},
}

var EmailRateLimitMiddleware = middleware.RateLimiterWithConfig(emailRateLimitConfig)
