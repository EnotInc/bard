package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultConfigFile = ".bard/config.json"
const configDir = ".bard"
const defaultTabStop = 4

// NOTE: is it alright to store config like that?
var global *Config

func Get() *Config {
	return global
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
	config := getCongfigPath()

	// creating a default config if bard.json is not found
	if _, err := os.Stat(config); err != nil {
		json, _ := json.MarshalIndent(defaultConfing, "", "    ")
		dir := getConfigDir()
		os.Mkdir(dir, 0755)
		os.WriteFile(config, []byte(json), 0644)
		global = defaultConfing
		return
	}

	data, err := os.ReadFile(config)
	if err != nil {
		global = defaultConfing
		return
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		global = defaultConfing
		return
	}

	global = cfg
	FixConfig()
}

func FixConfig() {
	if global.TabStop <= 0 {
		global.TabStop = defaultTabStop
	}
}

// saving current configuration
func Save() {
	config := getCongfigPath()
	json, _ := json.MarshalIndent(global, "", "    ")
	os.WriteFile(config, []byte(json), 0644)
}

func getDefaultConfig() *Config {
	config := &Config{
		RLN:       false,
		ShowMD:    false,
		Render:    true,
		TabNames:  true,
		ThemeName: defaultThemeName,
		TabStop:   defaultTabStop,
	}
	return config
}
