package application

import (
	"github.com/labstack/echo/v4"
	"go-security/internal/service"
	"go-security/internal/web/controller"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ApplicationContext struct {
	Engine      *echo.Echo
	SqlEngine   *gorm.DB
	AppConfig   *Config
	Controllers []controller.Controller
	Models      []any
	Services    []service.IService
}

func (ctx *ApplicationContext) RegisterControllers(controllers ...controller.Controller) {
	for _, _controller := range controllers {
		ctx.Controllers = append(ctx.Controllers, _controller)
	}
}
