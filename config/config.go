package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getCongfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "bard"
	}
	return filepath.Join(home, "bard")
}

func InitConfig() *Config {
	defaultConfing := getDefaultConfig()
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
	json, _ := json.MarshalIndent(c, "", "    ")
	os.WriteFile(path, []byte(json), 0644)
}

func getDefaultConfig() *Config {
	config := &Config{
		Editor: Editor{
			RLN:      false,
			ShowMD:   false,
			Render:   true,
			TabNames: true,
		},
		Theme: Theme{
			General: General{
				LineNumber:  "\033[90m",
				CurrentLine: "\033[33m",
				BottomBar:   "\033[48;5;16m",
				Selection:   "\033[100m",
				Command:     "\033[33m",
				EmptyLine:   "\033[36m",
				Tab:         "\033[94m",
			},
			Markdown: Markdown{
				Header1:    "\033[34m",
				Header2:    "\033[34m",
				Header3:    "\033[34m",
				Header4:    "\033[34m",
				Header5:    "\033[34m",
				Header6:    "\033[34m",
				Highlight:  "\033[43m",
				Symbol:     "\033[90m",
				Quote:      "\033[32m",
				NumberList: "\033[35m",
				Tag:        "\033[35m",
				CodeBg:     "\033[48;5;234m",
				CodeText:   "\033[33m",
				Image:      "\033[4;36m",
				Link:       "\033[4;36m",
			},
			Code: Code{
				Background: "\033[48;5;234m",
				Keyword:    "\033[33m",
				String:     "\033[32m",
				Number:     "\033[35m",
				Bracket:    "\033[35m",
				Symbol:     "\033[33m",
				Comment:    "\033[90m",
			},
		},
	}
	return config
}
