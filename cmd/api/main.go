package main

import (
	"go-security/security/application"
)

func main() {

	config := application.MustNewAppConfigFromFile("/Users/william_w_chen/Desktop/tofu-erp/cmd/api/config.yaml")
	app := application.MustNewApplication(config)
	context := application.MustNewSecurityApplicationContext(app)
	app.InjectContextCollection(context)
	app.Run()
}
