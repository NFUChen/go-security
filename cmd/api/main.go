package main

import (
	"go-security/internal/application"
)

func main() {
	config := application.MustNewAppConfigFromFile("./config.yaml")
	app := application.MustNewApplication(config)
	context := application.MustNewSecurityApplicationContext(config, app.SqlEngine, app.Engine)
	app.InjectContextCollection(context)
	app.Run()
}
