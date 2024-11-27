package main

import (
	erpApplication "go-security/erp/application"
	"go-security/security/application"
)

func main() {

	config := application.MustNewAppConfigFromFile("/Users/william_w_chen/Desktop/tofu-erp/cmd/api/config.yaml")
	app := application.MustNewApplication(config)
	context := application.MustNewSecurityApplicationContext(config, app.SqlEngine, app.Engine)
	erpAppConfig := erpApplication.MustNewErpApplicationConfig("/Users/william_w_chen/Desktop/tofu-erp/cmd/api/erp.yaml")
	erpContext := erpApplication.MustNewErpApplicationContext(erpAppConfig, app, context)
	app.InjectContextCollection(context, erpContext)
	app.Run()
}
