package controller

import (
	_ "github.com/joho/godotenv/autoload"
)

type Controller interface {
	RegisterRoutes()
}
