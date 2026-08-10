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
	theme := &Theme{
		General: General{
			SelectedTile: "#0000ff",
			LineNumber:   "#555555",
			CurrentLine:  "#ffff00",
			BottomBar:    "#000000",
			Selection:    "#cccccc",
			Command:      "#ffff00",
			EmptyLine:    "#00ffff",
			Message:      "#00ff00",
			Error:        "#ff0000",
			Tab:          "#0000ff",
		},
		Markdown: Markdown{
			Header1:    "#0000ff",
			Header2:    "#0000ff",
			Header3:    "#0000ff",
			Header4:    "#0000ff",
			Header5:    "#0000ff",
			Header6:    "#0000ff",
			Highlight:  "#ffff00",
			HTMLSymbol: "#555555",
			HTMLText:   "#ff0000",
			Symbol:     "#555555",
			Quote:      "#00ff00",
			NumberList: "#ff00ff",
			Tag:        "#ff00ff",
			CodeLineBg: "#1c1c1c",
			CodeHeader: "#1c1c1c",
			CodeText:   "#ffff00",
			Image:      "#00ffff",
			Link:       "#00ffff",
		},
		Code: Code{
			Background: "#1c1c1c",
			Keyword:    "#ffff00",
			String:     "#00ff00",
			Number:     "#ff00ff",
			Bracket:    "#ff00ff",
			Symbol:     "#ffff00",
			Comment:    "#555555",
		},
	}
	return theme
}
