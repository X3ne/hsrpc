package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/X3ne/hsrpc/src/consts"
	"github.com/X3ne/hsrpc/src/utils"
)

func LoadConfig() (AppConfig, error) {
	appData, err := utils.GetAppPath()
	if err != nil {
		return AppConfig{}, err
	}
	appDataPath := filepath.Join(appData, consts.ConfigFile)

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
	appData, err := utils.GetAppPath()
	if err != nil {
		return err
	}
	appDataPath := filepath.Join(appData, consts.ConfigFile)
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = utils.CreateDir(appDataPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(appDataPath, configJSON, 0600)
	if err != nil {
		return err
	}

	return nil
}
