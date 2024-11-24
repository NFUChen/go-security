package controller

import (
	_ "github.com/joho/godotenv/autoload"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Controller interface {
	RegisterRoutes()
}
