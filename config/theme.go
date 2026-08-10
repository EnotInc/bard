package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var theme *Theme

func setTheme(t *Theme) {
	t.parceColors()
	theme = t
}

func GetTheme() *Theme {
	return theme
}

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

type apply bool

const (
	foreground apply = true
	background apply = false
)

func hexToAscii(c string, a apply) (string, error) {
	hex := strings.TrimPrefix(c, "#")
	if len(hex) != 6 {
		return "", fmt.Errorf("Invalid hex len in string: %s", c)
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%2x%2x%2x", &r, &g, &b)
	if err != nil {
		return "", err
	}

	var escape string
	switch a {
	case foreground:
		escape = fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	case background:
		escape = fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	}

	return escape, nil
}

func (t *Theme) parceColors() {
	t.parceRecursive(reflect.ValueOf(t).Elem())
}

func (t *Theme) parceRecursive(val reflect.Value) {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Kind() == reflect.Struct {
			t.parceRecursive(field)
			continue
		}

		if field.CanSet() && field.Kind() == reflect.String {
			hex := field.String()
			apply := fieldType.Tag.Get("type")

			var ascii string
			var err error

			switch apply {
			case "foreground":
				ascii, err = hexToAscii(hex, foreground)
			case "background":
				ascii, err = hexToAscii(hex, background)
			default:
				continue
			}

			if err != nil {
				panic(err)
			}

			field.SetString(ascii)
		}
	}
}

func InitTheme(themeName string) error {
	defaultTheme := getDefaultTheme()
	theme_path := getThemePath(themeName)

	if _, err := os.Stat(theme_path); err != nil {
		json, _ := json.MarshalIndent(defaultTheme, "", "    ")
		dir := getThemeDir()
		os.Mkdir(dir, 0755)
		os.WriteFile(getThemePath(defaultThemeName), []byte(json), 0644)
		setTheme(defaultTheme)
		return fmt.Errorf("Unable to load theme %s", themeName)
	}

	data, err := os.ReadFile(theme_path)
	if err != nil {
		theme = defaultTheme
		return fmt.Errorf("Unable to load theme %s", themeName)
	}

	t := &Theme{}
	err = json.Unmarshal(data, t)
	if err != nil {
		setTheme(defaultTheme)
		return fmt.Errorf("Unable to load theme %s", themeName)
	}

	setTheme(t)
	return nil
}

func ChangeTheme(themeName string) string {
	theme_path := getThemePath(themeName)

	if _, err := os.Stat(theme_path); err != nil {
		return "Unknown theme '" + themeName + "'"
	}

	data, err := os.ReadFile(theme_path)
	if err != nil {
		return "unable to read theme file '" + themeName + "'"
	}

	new := &Theme{}
	err = json.Unmarshal(data, new)
	if err != nil {
		return "unable to set theme '" + themeName + "'"
	}

	setTheme(new)
	return ""
}

func getDefaultTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#8c9dff",
			LineNumber:   "#9aa0a6",
			CurrentLine:  "#f5c542",
			BottomBar:    "#2d2d3f",
			Selection:    "#d4d4d4",
			Command:      "#f5c542",
			EmptyLine:    "#6fc1d4",
			Message:      "#7fc97f",
			Error:        "#e8747c",
			Tab:          "#7c9dff",
		},
		Markdown: Markdown{
			Header1:    "#7c9dff",
			Header2:    "#7c9dee",
			Header3:    "#7c9ddd",
			Header4:    "#7c9dcc",
			Header5:    "#7c9dbb",
			Header6:    "#7c9daa",
			Highlight:  "#f5e68b",
			HTMLSymbol: "#9aa0a6",
			HTMLText:   "#e8747c",
			Symbol:     "#9aa0a6",
			Quote:      "#7fc97f",
			NumberList: "#c997d9",
			Tag:        "#c997d9",
			CodeLineBg: "#2a2a3a",
			CodeHeader: "#2a2a3a",
			CodeText:   "#f5c542",
			Image:      "#6fc1d4",
			Link:       "#6fc1d4",
		},
		Code: Code{
			Background: "#2a2a3a",
			Keyword:    "#f5c542",
			String:     "#7fc97f",
			Number:     "#c997d9",
			Bracket:    "#c997d9",
			Symbol:     "#f5c542",
			Comment:    "#9aa0a6",
		},
	}
}
