package security

import (
	"gopkg.in/yaml.v3"
	"os"
)

func MustNewConfigFromFile[T any](configPath string) *T {
	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var config T
	if err := yaml.Unmarshal(file, &config); err != nil {
		panic(err)
	}
	return &config
}
