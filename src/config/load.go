package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/X3ne/hsrpc/src/consts"
	"github.com/X3ne/hsrpc/src/utils"
)

func LoadConfig() (AppConfig, error) {
	appData, err := utils.GetAppPath()
	if err != nil {
		return AppConfig{}, err
	}
	appDataPath := filepath.Join(appData, consts.ConfigFile)

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

	defaultConfig := NewConfig()
	updateConfigWithDefaults(&config, defaultConfig)

	err = SaveConfig(config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

func updateConfigWithDefaults(config *AppConfig, defaultConfig AppConfig) {
	configType := reflect.TypeOf(*config)

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)

		fieldValue := reflect.ValueOf(config).Elem().FieldByName(field.Name)

		defaultValue := reflect.ValueOf(defaultConfig).FieldByName(field.Name)

		if fieldValue.Kind() == reflect.Struct {
			for j := 0; j < fieldValue.NumField(); j++ {
				nestedField := fieldValue.Field(j)
				if nestedField.IsZero() {
					nestedField.Set(defaultValue.Field(j))
				}
			}
		} else {
			if fieldValue.IsZero() {
				if fieldValue.Kind() != reflect.Bool {
					fieldValue.Set(defaultValue)
				}
			}
		}
	}
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
