package application

import "go-security/security"

type AwsConfig struct {
	Region             string `yaml:"region"`
	AwsAccessKeyID     string `yaml:"access_key_id"`
	AwsSecretAccessKey string `yaml:"secret_access_key"`
}

type LineConfig struct {
	ChannelSecret      string `yaml:"channel_secret"`
	ChannelAccessToken string `yaml:"channel_access_token"`
}

type ErpApplicationConfig struct {
	Aws  *AwsConfig  `yaml:"aws"`
	Line *LineConfig `yaml:"line"`
}

func MustNewErpApplicationConfig(configPath string) *ErpApplicationConfig {
	return security.MustNewConfigFromFile[ErpApplicationConfig](configPath)
}
