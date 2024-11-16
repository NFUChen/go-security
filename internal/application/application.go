package application

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-security/internal/repository"
	"go-security/internal/service"
	"go-security/internal/web/controller"
	web "go-security/internal/web/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"path/filepath"
	"strconv"
)

type Application struct {
	Engine      *echo.Echo
	AppConfig   *Config
	SqlEngine   *gorm.DB
	Controllers []controller.Controller
}

func (app *Application) migrateDatabase() {
	err := app.SqlEngine.AutoMigrate(repository.GetAllModels()...)
	if err != nil {
		panic(err)
	}
}

func setupLogger() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
}

func MustNewApplication(config *Config) *Application {
	setupLogger()
	engine := echo.New()

	sqlEngine, err := gorm.Open(postgres.Open(config.PostgresDataSource.AsDSN()), &gorm.Config{})
	log.Printf("Connected to database: %s", config.PostgresDataSource.DatabaseName)

	userRepo := repository.NewUserRepository(sqlEngine)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userService, config.Security.Secret)
	authMiddleware := web.NewAuthMiddleware(authService, config.Security.ExcludedRoutePrefixes)
	log.Printf("Security excluded routes: %v", config.Security.ExcludedRoutePrefixes)

	baseRouterGroup := engine.Group("/api")

	mainController := controller.NewMainController(engine)
	authController := controller.NewAuthController(baseRouterGroup, authService, userService)
	userController := controller.NewUserController(baseRouterGroup, userService)
	controllers := []controller.Controller{
		mainController,
		authController,
		userController,
	}

	if err != nil {
		panic(err)
	}

	engine.Use(middleware.Recover())
	engine.Use(middleware.Logger())
	engine.Use(web.ErrorMiddlewareFunc)
	engine.Use(authMiddleware.AuthMiddlewareFunc)
	engine.Use(web.CORSMiddlewareFunc)

	return &Application{
		AppConfig:   config,
		Engine:      engine,
		Controllers: controllers,
		SqlEngine:   sqlEngine,
	}
}

func (app *Application) RegisterControllers() {
	for _, _controller := range app.Controllers {
		_controller.RegisterRoutes()
	}
}

func (app *Application) Run() {
	app.migrateDatabase()
	app.RegisterControllers()
	err := app.Engine.Start(fmt.Sprintf(":%d", app.AppConfig.Server.Port))
	if err != nil {
		panic(err)
	}

}
