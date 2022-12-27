package configuration

import (
	"chimera/utility/logging"
	"os"

	"gopkg.in/yaml.v3"
)

func GetConfiguration() ConfigData {
	cfg := ConfigData{}

	yamlFile, err := os.ReadFile("cfg/next.yaml")
	if err != nil {
		logging.Error("YAML Configuration", "Failed to Read Configuration File! (%s)", err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		logging.Error("YAML Configuration", "Failed to Unmarshal Configuration File! (%s)", err.Error())
	}

	return cfg
}
