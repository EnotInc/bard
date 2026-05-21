package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const defaultThemeName = "bard.json"
const themeDir = ".bard/themes"

func (c *Config) DefaultThemeName() string {
	return defaultThemeName
}

func getThemeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return themeDir
	}
	return filepath.Join(home, themeDir)
}

func getThemePath(themeName string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(themeDir, themeName)
	}
	return filepath.Join(home, themeDir, themeName)
}

func InitTheme(themeName string) (*Theme, error) {
	defaultTheme := getDefaultTheme()
	theme := getThemePath(themeName)

	if _, err := os.Stat(theme); err != nil {
		json, _ := json.MarshalIndent(defaultTheme, "", "    ")
		dir := getThemeDir()
		os.Mkdir(dir, 0755)
		os.WriteFile(getThemePath(defaultThemeName), []byte(json), 0644)
		return defaultTheme, fmt.Errorf("Unable to load theme %s", themeName)
	}

	data, err := os.ReadFile(theme)
	if err != nil {
		return defaultTheme, fmt.Errorf("Unable to load theme %s", themeName)
	}

	t := &Theme{}
	err = json.Unmarshal(data, t)
	if err != nil {
		return defaultTheme, fmt.Errorf("Unable to load theme %s", themeName)
	}

	return t, nil
}

func (t *Theme) ChangeTheme(themeName string) string {
	theme := getThemePath(themeName)

	if _, err := os.Stat(theme); err != nil {
		return "Unknown theme '" + themeName + "'"
	}

	data, err := os.ReadFile(theme)
	if err != nil {
		return "unable to read theme file '" + themeName + "'"
	}

	new := &Theme{}
	err = json.Unmarshal(data, t)
	if err != nil {
		return "unable to set theme '" + themeName + "'"
	}

	t = new // just in case
	return ""
}

func getDefaultTheme() *Theme {
	theme := &Theme{
		General: General{
			LineNumber:  "\033[90m",
			CurrentLine: "\033[33m",
			BottomBar:   "\033[48;5;16m",
			Selection:   "\033[100m",
			Command:     "\033[33m",
			EmptyLine:   "\033[36m",
			Message:     "\033[31m",
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
			HTMLSymbol: "\033[90m",
			HTMLText:   "\033[31m",
			Symbol:     "\033[90m",
			Quote:      "\033[32m",
			NumberList: "\033[35m",
			Tag:        "\033[35m",
			CodeLineBg: "\033[48;5;234m",
			CodeHeader: "\033[48;5;234m",
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
	}
	return theme
}
