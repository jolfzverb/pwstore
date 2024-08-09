package dependencies

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigType struct {
	Database struct {
		ConnectionString string `yaml:"connectionString"`
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
