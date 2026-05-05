package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultConfigFile = ".bard/config.json"
const configDir = ".bard"

func getCongfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return defaultConfigFile
	}
	return filepath.Join(home, defaultConfigFile)
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return configDir
	}
	return filepath.Join(home, configDir)
}

func InitConfig() *Config {
	defaultConfing := getDefaultConfig()
	config := getCongfigPath()

	// creating a default config if bard.json is not found
	if _, err := os.Stat(config); err != nil {
		json, _ := json.MarshalIndent(defaultConfing, "", "    ")
		dir := getConfigDir()
		os.Mkdir(dir, 0755)
		os.WriteFile(config, []byte(json), 0644)
		return defaultConfing
	}

	data, err := os.ReadFile(config)
	if err != nil {
		return defaultConfing
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return defaultConfing
	}

	return cfg
}

// saving current configuration
func (c *Config) Save() {
	config := getCongfigPath()
	json, _ := json.MarshalIndent(c, "", "    ")
	os.WriteFile(config, []byte(json), 0644)
}

func getDefaultConfig() *Config {
	config := &Config{
		RLN:       false,
		ShowMD:    false,
		Render:    true,
		TabNames:  true,
		ThemeName: defaultThemeName,
	}
	return config
}
