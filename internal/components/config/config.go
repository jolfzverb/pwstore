package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Model struct {
	OpenIDSettings struct {
		AuthorizationEndpoint string   `yaml:"authorizationEndpoint"`
		ClientID              string   `yaml:"clientId"`
		RedirectURI           string   `yaml:"redirectUri"`
		ResponseType          string   `yaml:"responseType"`
		Scope                 []string `yaml:"scope"`
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
