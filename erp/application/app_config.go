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

type ErpApplicationConfig struct {
	Aws   *AwsConfig           `yaml:"aws"`
	Line  *service.LineConfig  `yaml:"line"`
	Minio *service.MinioConfig `yaml:"minio"`
}

func MustNewErpApplicationConfig(configPath string) *ErpApplicationConfig {
	return security.MustNewConfigFromFile[ErpApplicationConfig](configPath)
}
