package application

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"path/filepath"
	"reflect"
	"strconv"
)

type Application struct {
	ContextCollection []*ApplicationContext
	AppConfig         *Config
	Engine            *echo.Echo
	SqlEngine         *gorm.DB
}

type Runnable interface {
	Run()
}

func (app *Application) migrateDatabase() {
	for _, _context := range app.ContextCollection {
		err := app.SqlEngine.AutoMigrate(_context.Models...)
		if err != nil {
			panic(err)
		}
	}

}

func (app *Application) setupLogger() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
}

func (app *Application) InjectContextCollection(appContextCollection ...*ApplicationContext) {
	app.ContextCollection = appContextCollection
}

func MustNewApplication(config *Config) *Application {
	engine := echo.New()
	sqlEngine, err := gorm.Open(postgres.Open(config.PostgresDataSource.AsDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Application{
		AppConfig: config,
		Engine:    engine,
		SqlEngine: sqlEngine,
	}
}

func (app *Application) registerControllerRoutes() {
	for _, _context := range app.ContextCollection {
		for _, _controller := range _context.Controllers {
			log.Info().Msgf("Registering routes for web: %s", reflect.TypeOf(_controller).String())
			_controller.RegisterRoutes()
		}
	}

}

func (app *Application) postConstructServices() {
	for _, _context := range app.ContextCollection {
		for _, _service := range _context.Services {
			log.Info().Msgf("PostConstruct for service: %s", reflect.TypeOf(_service).String())
			_service.PostConstruct()
		}
	}
	log.Info().Msg("PostConstruct for all services completed")
}

func (app *Application) Run() {
	log.Printf("Starting application...")
	fmt.Printf("%s\n", app.AppConfig.AsJson())
	app.migrateDatabase()
	app.postConstructServices()
	app.registerControllerRoutes()
	err := app.Engine.Start(fmt.Sprintf(":%d", app.AppConfig.Server.Port))
	if err != nil {
		panic(err)
	}

}
