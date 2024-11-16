package application

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type SecurityConfig struct {
	Secret                string   `yaml:"secret"`
	ExcludedRoutePrefixes []string `yaml:"excluded_routes_prefixes"`
}

type PostgresDataSourceConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"db_name"`
}

func (config *PostgresDataSourceConfig) AsDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName,
	)
}

type Config struct {
	Server             ServerConfig             `yaml:"server"`
	Security           SecurityConfig           `yaml:"security"`
	PostgresDataSource PostgresDataSourceConfig `yaml:"postgres_data_source"`
}

func MustNewConfigFromFile(configPath string) *Config {
	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		panic(err)
	}
	return &config
}
