package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const defaultThemeName = "bard.json"
const defaultConfigFile = ".bard/config.json"
const configDir = ".bard"
const defaultTabStop = 4
const defaultResizeDuration = 200

var config *Config

func GetConfig() *Config {
	return config
}

func (c *Config) DefaultThemeName() string {
	return defaultThemeName
}

func getConfigPath() string {
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

func ReadConfig() error {
	defaultConfing := getDefaultConfig()
	cfg_path := getConfigPath()

	data, err := os.ReadFile(cfg_path)
	if err != nil {
		config = defaultConfing
		return fmt.Errorf("unable to read config file. Using default config instead")
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		config = defaultConfing
		return fmt.Errorf("Unable to parce config. Using default settings instead")
	}

	config = cfg
	FixConfig()
	return nil
}

func InitConfig() {
	defaultConfing := getDefaultConfig()
	cfg_path := getConfigPath()

	// creating a default config if bard.json is not found
	if _, err := os.Stat(cfg_path); err != nil {
		json, _ := json.MarshalIndent(defaultConfing, "", "    ")
		dir := getConfigDir()
		os.Mkdir(dir, 0755)
		os.WriteFile(cfg_path, []byte(json), 0644)
		config = defaultConfing
		return
	}

	ReadConfig()
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
	cfg := getConfigPath()
	json, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(cfg, []byte(json), 0644)
	if err != nil {
		panic(err)
	}
}

func getDefaultConfig() *Config {
	return &Config{
		RLN:        false,
		ShowMD:     false,
		ShowEmpty:  true,
		Render:     true,
		TabNames:   true,
		ThemeName:  defaultThemeName,
		TabStop:    defaultTabStop,
		ResizeTime: defaultResizeDuration,
		KeepTabs:   true,
		ShowIcons:  true,
		ShowBorder: false,
		ShowDot:    true,
	}
}
