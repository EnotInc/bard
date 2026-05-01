package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getCongfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".bard/config.json"
	}
	return filepath.Join(home, ".bard/config.json")
}

func InitConfig() *Config {
	defaultConfing := getDefaultConfig()
	config := getCongfigPath()

	// creating a default config if bard.json is not found
	if _, err := os.Stat(config); err != nil {
		json, _ := json.MarshalIndent(defaultConfing, "", "    ")
		dir := getDirPath()
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
		RLN:      false,
		ShowMD:   false,
		Render:   true,
		TabNames: true,
	}
	return config
}
