package main

import (
	"go-security/internal/application"
)

func main() {
	config := application.MustNewConfigFromFile("./config.yaml")
	app := application.MustNewApplication(config)
	app.Run()
}
