package main

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	Tag1 string `json:"tag1"`
	Tag2 string `json:"tag2"`
	Tag3 string `json:"tag3"`
}

const configPath = "config.json"

// LoadConfig reads the tags from the local file
func LoadConfig() AppConfig {
	file, err := os.ReadFile(configPath)
	if err != nil {
		// Return defaults if file doesn't exist
		return AppConfig{Tag1: "Default", Tag2: "Default", Tag3: "Default"}
	}

	var config AppConfig
	json.Unmarshal(file, &config)
	return config
}

// SaveConfig writes the tags to the local file
func SaveConfig(config AppConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
