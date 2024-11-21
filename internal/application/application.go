package application

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strconv"
)

type Application struct {
	Context *ApplicationContext
}

type Runnable interface {
	Run()
}

func (app *Application) migrateDatabase() {
	err := app.Context.SqlEngine.AutoMigrate(app.Context.Models...)
	if err != nil {
		panic(err)
	}
}

func (app *Application) setupLogger() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
}

func MustNewApplication(appContext *ApplicationContext) *Application {
	app := &Application{
		Context: appContext,
	}
	app.setupLogger()
	return app
}

func (app *Application) registerControllerRoutes() {
	for _, _controller := range app.Context.Controllers {
		_controller.RegisterRoutes()
	}
}

func (app *Application) postConstructServices() {
	for _, _service := range app.Context.Services {
		_service.PostConstruct()
	}
}

func (app *Application) Run() {
	log.Printf("Starting application...")
	fmt.Printf("%s\n", app.Context.AppConfig.AsJson())
	app.migrateDatabase()
	app.postConstructServices()
	app.registerControllerRoutes()
	err := app.Context.Engine.Start(fmt.Sprintf(":%d", app.Context.AppConfig.Server.Port))
	if err != nil {
		panic(err)
	}

}
