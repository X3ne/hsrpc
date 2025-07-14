package main

type ScriptConfig struct {
	OutputPath string
	AssetsPath string
}

func InitConfig() *ScriptConfig {
	return &ScriptConfig{
		OutputPath: "../src-tauri/game-data/",
		AssetsPath: "assets/",
	}
}
