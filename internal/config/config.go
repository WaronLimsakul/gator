package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory")
	}
	const configName string = ".gatorconfig.json"
	return homeDir + "/" + configName, nil
}

func writeConfig(conf Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file path")
	}

	bytesConf, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("error encoding config struct")
	}

	if err = os.WriteFile(configFilePath, bytesConf, 0600); err != nil {
		return fmt.Errorf("error writing a config file")
	}

	return nil
}

func ReadConfig() (Config, error) {
	// note that the path doesn't have "/"
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting a config file path")
	}

	configStream, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("error accessing a config file")
	}

	var conf Config
	if err = json.Unmarshal(configStream, &conf); err != nil {
		return Config{}, fmt.Errorf("error decoding a config file")
	}
	return conf, nil
}

func (c *Config) SetUser(userName string) error {
	conf, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("error reading a config file")
	}

	conf.CurrentUsername = userName
	err = writeConfig(conf)
	if err != nil {
		return fmt.Errorf("error writing a config file")
	}
	return nil
}
