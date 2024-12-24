package application

import (
	"go-security/erp/internal/service"
	"go-security/security"
)

type AwsConfig struct {
	Region             string `yaml:"region"`
	AwsAccessKeyID     string `yaml:"access_key_id"`
	AwsSecretAccessKey string `yaml:"secret_access_key"`
}

type RedisConfig struct {
	Address       string `yaml:"address"`
	Password      string `yaml:"password"`
	DataBaseIndex int    `yaml:"database_index"`
}

type ErpApplicationConfig struct {
	Aws   *AwsConfig           `yaml:"aws"`
	Line  *service.LineConfig  `yaml:"line"`
	Minio *service.MinioConfig `yaml:"minio"`
	Redis *RedisConfig         `yaml:"redis"`
}

func MustNewErpApplicationConfig(configPath string) *ErpApplicationConfig {
	return security.MustNewConfigFromFile[ErpApplicationConfig](configPath)
}
