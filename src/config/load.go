package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	appDataDir = "hsrpc"
	configFile = "config.txt"
)

func LoadConfig() (AppConfig, error) {
	appDataPath := filepath.Join(os.Getenv("APPDATA"), appDataDir, configFile)

	// Config file does not exist, create a new one with default values
	if _, err := os.Stat(appDataPath); os.IsNotExist(err) {
		defaultConfig := NewConfig()
		err := SaveConfig(defaultConfig)
		if err != nil {
			return AppConfig{}, err
		}
		return defaultConfig, nil
	}

	fileContent, err := os.ReadFile(appDataPath)
	if err != nil {
		return AppConfig{}, err
	}

	var config AppConfig
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

func SaveConfig(config AppConfig) error {
	appDataPath := filepath.Join(os.Getenv("APPDATA"), appDataDir, configFile)
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(appDataPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(appDataPath, configJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
