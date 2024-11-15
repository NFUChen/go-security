package application

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-security/internal/web"
)

type Application struct {
	Engine      *echo.Echo
	AppConfig   *Config
	Controllers []web.Controller
}

func NewApplication(config *Config) *Application {
	engine := echo.New()
	mainController := web.NewMainController(engine)
	controllers := []web.Controller{
		mainController,
	}

	return &Application{
		AppConfig:   config,
		Engine:      engine,
		Controllers: controllers,
	}
}

func (app *Application) RegisterControllers() {
	for _, controller := range app.Controllers {
		controller.RegisterRoutes()
	}
}

func (app *Application) Run() {
	app.RegisterControllers()
	err := app.Engine.Start(fmt.Sprintf(":%d", app.AppConfig.Server.Port))
	if err != nil {
		panic(err)
	}

}
