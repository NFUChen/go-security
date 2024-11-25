package application

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"go-security/security/repository"
	"go-security/security/service"
	"go-security/security/service/oauth"
	"go-security/security/web/controller"
	web "go-security/security/web/middleware"
	"gorm.io/gorm"
)

func MustNewSecurityApplicationContext(config *Config, sqlEngine *gorm.DB, engine *echo.Echo) *ApplicationContext {
	log.Printf("Starting application...")
	fmt.Printf("%s\n", config.AsJson())
	ctx := context.Background()
	log.Info().Msgf("Connected to database: %s", config.PostgresDataSource.DatabaseName)

	otpService := service.NewOtpService(service.GenerateOtpCode)
	userRepo := repository.NewUserRepository(sqlEngine)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userService, config.Security.Secret)
	authMiddleware := web.NewAuthMiddleware(authService, config.Security.ExcludedRoutePrefixes)
	log.Info().Msgf("Security excluded routes: %v", config.Security.ExcludedRoutePrefixes)
	smtpService := service.NewSmtpService(ctx, &config.Smtp)
	resetPasswordService := service.NewUserResetPasswordService(smtpService, userService, authService, otpService)
	verificationService := service.NewUserVerificationService(smtpService, userService, authService, otpService)

	googleAuthService := oauth.NewGoogleAuthService(&config.GoogleAuthConfig, authService, userService)

	baseRouterGroup := engine.Group("/api")

	mainController := controller.NewMainController(engine)
	authController := controller.NewAuthController(baseRouterGroup, authService, userService, verificationService, resetPasswordService)
	userController := controller.NewUserController(baseRouterGroup, userService, resetPasswordService, verificationService)
	googleAuthController := controller.NewGoogleAuthController(baseRouterGroup, googleAuthService, config.Security.AuthRedirectUrl)
	controllers := []controller.Controller{
		mainController,
		authController,
		userController,
		googleAuthController,
	}
	middlewares := []echo.MiddlewareFunc{
		middleware.Recover(),
		middleware.Logger(),
		web.ErrorMiddlewareFunc,
		authMiddleware.AuthMiddlewareFunc,
		web.CORSMiddlewareFunc,
	}

	engine.Use(middlewares...)

	services := []service.IService{
		userService,
		authService,
		smtpService,
	}

	appContext := &ApplicationContext{
		Controllers: controllers,
		Models:      repository.NewSecurityModelProvider().ProvideModels(),
		Services:    services,
	}
	return appContext
}
