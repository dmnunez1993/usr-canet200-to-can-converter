package usrcanettocan

import (
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "./config.yaml"

type UsrCanetDeviceConverterConfig struct {
	Host   string `yaml:"host" json:"host"`
	Port   int64  `yaml:"port" json:"port"`
	Target string `yaml:"target" json:"target"`
}

type UsrCanetConverterConfig struct {
	DeviceConverters []UsrCanetDeviceConverterConfig `yaml:"device_converters" json:"deviceConverters"`
}

func getConfigPath() string {
	path, found := os.LookupEnv("CONFIG_PATH")

	if found {
		return path
	}

	return defaultConfigPath
}

func (d *UsrCanetConverterConfig) Save() error {
	source, err := yaml.Marshal(d)

	if err != nil {
		return err
	}

	err = os.WriteFile(getConfigPath(), source, 0644)

	if err != nil {
		return err
	}

	return nil
}

func GetConverterConfig() (UsrCanetConverterConfig, error) {
	yamlFile, err := os.ReadFile(getConfigPath())

	if err != nil {
		return UsrCanetConverterConfig{}, err
	}

	config := &UsrCanetConverterConfig{}
	err = yaml.Unmarshal(yamlFile, config)

	return *config, err
}

func GetOrCreateConverterConfig() (UsrCanetConverterConfig, error) {
	config, err := GetConverterConfig()

	if err != nil {
		config := UsrCanetConverterConfig{}
		err = config.Save()
	}

	return config, err
}
