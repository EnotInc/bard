package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getDirPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".bard/"
	}
	return filepath.Join(home, ".bard/")
}

func (c *Theme) Save() {
	theme := getCongfigPath()
	json, _ := json.MarshalIndent(theme, "", "    ")
	os.WriteFile(theme, []byte(json), 0644)
}

func getThemePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".bard/theme.json"
	}
	return filepath.Join(home, ".bard/theme.json")
}

func InitTheme() *Theme {
	defaultTheme := getDefaultTheme()
	theme := getThemePath()

	if _, err := os.Stat(theme); err != nil {
		json, _ := json.MarshalIndent(defaultTheme, "", "    ")
		dir := getDirPath()
		os.Mkdir(dir, 0755)
		os.WriteFile(theme, []byte(json), 0644)
		return defaultTheme
	}

	data, err := os.ReadFile(theme)
	if err != nil {
		return defaultTheme
	}

	t := &Theme{}
	err = json.Unmarshal(data, t)
	if err != nil {
		return defaultTheme
	}

	return t
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
	}
	return theme
}
