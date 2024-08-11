package secrets

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Model struct {
	Database struct {
		ConnectionString string `yaml:"connectionString"`
	} `yaml:"database"`
	OpenIDSettings struct {
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"openIdSettings"`
}

func GetConfig(filename string) (*Model, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config Model
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return &config, nil
}
