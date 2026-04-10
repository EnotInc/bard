package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// About |Config|
// |RLN| - relative line munber
// |ShowMD| - always show markdown symbols
// |Render| - enable or disable render
// |TabNames| - show file name is tabs or not
type Config struct {
	RLN       bool `json:"rln"`
	ShowMD    bool `json:"showmd"`
	Render    bool `json:"render"`
	TabNames  bool `json:"tab_names"`
	IsChanged bool
}

func getCongfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "bard.json"
	}
	return filepath.Join(home, "bard.json")
}

func InitConfig() *Config {
	defaultConfing := &Config{RLN: false, ShowMD: false, Render: true, TabNames: true, IsChanged: false}
	config := getCongfigPath()

	// creating a default config if bard.json is not found
	if _, err := os.Stat(config); err != nil {
		json, _ := json.Marshal(defaultConfing)
		os.WriteFile(config, []byte(json), 0644)
		return defaultConfing
	}

	data, err := os.ReadFile(config)
	if err != nil {
		return defaultConfing
	}
	cfg := &Config{}
	json.Unmarshal(data, cfg)
	return cfg
}

// saving current configuration
func (c *Config) Save() {
	path := getCongfigPath()
	json, _ := json.Marshal(c)
	os.WriteFile(path, []byte(json), 0644)
}
