package dependencies

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigType struct {
	Database struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"database"`
}

func GetConfig(filename string) (*ConfigType, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config ConfigType
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
