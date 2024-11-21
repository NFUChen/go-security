package main

import (
	"go-security/internal/application"
)

func main() {
	config := application.MustNewConfigFromFile("./config.yaml")
	context := application.MustNewSecurityApplicationContext(config)
	app := application.MustNewApplication(context)
	app.Run()
}
