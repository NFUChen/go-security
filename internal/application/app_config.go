package application

import (
	"encoding/json"
	"go-security/internal"
	"go-security/internal/repository"
	"go-security/internal/service"
	"go-security/internal/service/oauth"
	"go-security/internal/web/controller"
)

type Config struct {
	Server             controller.ServerConfig             `yaml:"server"`
	Security           service.SecurityConfig              `yaml:"security"`
	PostgresDataSource repository.PostgresDataSourceConfig `yaml:"postgres_data_source"`
	Smtp               service.SmtpConfig                  `yaml:"smtp"`
	GoogleAuthConfig   oauth.GoogleAuthConfig              `yaml:"google_auth"`
}

func (config *Config) AsJson() string {
	_json, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		return ""
	}
	return string(_json)
}

func MustNewAppConfigFromFile(configPath string) *Config {
	return internal.MustNewConfigFromFile[Config](configPath)
}
