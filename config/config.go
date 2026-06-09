package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultConfigFile = ".bard/config.json"
const configDir = ".bard"
const defaultTabStop = 4
const defaultResizeDuration = 200

// NOTE: is it alright to store config like that?
var config *Config

func GetConfig() *Config {
	return config
}

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

func InitConfig() {
	defaultConfing := getDefaultConfig()
	cfg_path := getCongfigPath()

	// creating a default config if bard.json is not found
	if _, err := os.Stat(cfg_path); err != nil {
		json, _ := json.MarshalIndent(defaultConfing, "", "    ")
		dir := getConfigDir()
		os.Mkdir(dir, 0755)
		os.WriteFile(cfg_path, []byte(json), 0644)
		config = defaultConfing
		return
	}

	data, err := os.ReadFile(cfg_path)
	if err != nil {
		config = defaultConfing
		return
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		config = defaultConfing
		return
	}

	config = cfg
	FixConfig()
}

func FixConfig() {
	if config.TabStop <= 0 {
		config.TabStop = defaultTabStop
	}

	if config.ResizeTime < 200 {
		config.ResizeTime = 200
	}
	if config.ResizeTime > 1000 {
		config.ResizeTime = 1000
	}
}

// saving current configuration
func Save() {
	cfg := getCongfigPath()
	json, _ := json.MarshalIndent(config, "", "    ")
	os.WriteFile(cfg, []byte(json), 0644)
}

func getDefaultConfig() *Config {
	config := &Config{
		RLN:        false,
		ShowMD:     false,
		Render:     true,
		TabNames:   true,
		ThemeName:  defaultThemeName,
		TabStop:    defaultTabStop,
		ResizeTime: defaultResizeDuration,
		KeepTabs:   true,
	}
	return config
}
