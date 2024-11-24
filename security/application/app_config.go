package application

import (
	"encoding/json"
	"go-security/security"
	"go-security/security/repository"
	"go-security/security/service"
	"go-security/security/service/oauth"
	"go-security/security/web/controller"
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
	return security.MustNewConfigFromFile[Config](configPath)
}
