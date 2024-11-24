package application

import (
	"go-security/security/service"
	"go-security/security/web/controller"
)

//goland:noinspection GoNameStartsWithPackageName
type ApplicationContext struct {
	Controllers []controller.Controller
	Models      []any
	Services    []service.IService
}
