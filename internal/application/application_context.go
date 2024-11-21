package application

import (
	"go-security/internal/service"
	"go-security/internal/web/controller"
)

//goland:noinspection GoNameStartsWithPackageName
type ApplicationContext struct {
	Controllers []controller.Controller
	Models      []any
	Services    []service.IService
}
