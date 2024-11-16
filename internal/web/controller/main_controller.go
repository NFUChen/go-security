package controller

import "github.com/labstack/echo/v4"

type MainController struct {
	Engine *echo.Echo
}

func NewMainController(engine *echo.Echo) *MainController {
	return &MainController{
		Engine: engine,
	}
}

func (controller MainController) RegisterRoutes() {
	controller.Engine.GET("/", func(context echo.Context) error {
		return context.String(200, "Hello, World!")
	})
}
