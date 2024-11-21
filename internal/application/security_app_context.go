package application

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"go-security/internal/repository"
	"go-security/internal/service"
	"go-security/internal/service/oauth"
	"go-security/internal/web/controller"
	web "go-security/internal/web/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustNewSecurityApplicationContext(config *Config) *ApplicationContext {
	log.Printf("Starting application...")
	fmt.Printf("%s\n", config.AsJson())
	ctx := context.Background()
	engine := echo.New()

	sqlEngine, err := gorm.Open(postgres.Open(config.PostgresDataSource.AsDSN()), &gorm.Config{})
	log.Info().Msgf("Connected to database: %s", config.PostgresDataSource.DatabaseName)

	otpService := service.NewOtpService()
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
	googleAuthController := controller.NewGoogleAuthController(baseRouterGroup, googleAuthService)
	controllers := []controller.Controller{
		mainController,
		authController,
		userController,
		googleAuthController,
	}

	if err != nil {
		panic(err)
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
		authService,
	}

	appContext := &ApplicationContext{
		Engine:    engine,
		SqlEngine: sqlEngine,

		AppConfig:   config,
		Controllers: controllers,
		Models:      repository.NewSecurityModelProvider().ProvideModels(),
		Services:    services,
	}
	return appContext
}
