package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/EnotInc/Bard/internal/services"
)

var theme *Theme

func setTheme(t *Theme) {
	t.parceColors()
	theme = t
}

func GetTheme() *Theme {
	return theme
}

const themeDir = ".bard/themes"

func GetThemeList() (map[string][9]string, error) {
	var path string
	home, err := os.UserHomeDir()
	if err != nil {
		path = themeDir
	}
	path = filepath.Join(home, themeDir)

	files, err := os.ReadDir(path)
	if err != nil {
		return map[string][9]string{}, err
	}

	themes := make(map[string][9]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			pallete, err := getThemePallete(file.Name())
			if err != nil {
				continue
			}
			themes[file.Name()] = pallete
		}
	}

	return themes, nil
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
				ascii, err = services.HexToAscii(hex, services.Foreground)
			case "background":
				ascii, err = services.HexToAscii(hex, services.Background)
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

func createAllThemes() {
	on_err := func(name string, theme *Theme) {
		json, _ := json.MarshalIndent(theme, "", "    ")
		dir := getThemeDir()
		os.Mkdir(dir, 0755)
		os.WriteFile(getThemePath(name), []byte(json), 0644)
		setTheme(theme)
	}

	themes := map[theme_name]*Theme{
		bard:        getDefaultTheme(),
		catppuccin:  getCatppuccinTheme(),
		dracula:     getDraculaTheme(),
		gruvbox:     getGruvboxTheme(),
		solarized:   getSolarizedTheme(),
		tokyo_night: getTokyoNightTheme(),
		monokai:     getMonokaiTheme(),
		one_dark:    getOneDarkTheme(),
		nord:        getNordTheme(),
	}

	for name, theme := range themes {
		path := getThemePath(string(name))

		if _, err := os.Stat(path); err != nil {
			on_err(string(name), theme)
		}
	}

}

func InitTheme(themeName string) error {
	defaultTheme := getDefaultTheme()
	theme_path := getThemePath(themeName)
	createAllThemes()

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

func getThemePallete(name string) ([9]string, error) {
	path := getThemePath(name)
	if _, err := os.Stat(path); err != nil {
		return [9]string{}, fmt.Errorf("Unknown theme: %s", name)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return [9]string{}, fmt.Errorf("Unable to read theme file '%s'", name)
	}

	tmp := &Theme{}
	err = json.Unmarshal(data, tmp)
	if err != nil {
		return [9]string{}, fmt.Errorf("Unable to parce theme '%s'", name)
	}
	tmp.parceColors()

	var pallete [9]string
	pallete[0] = tmp.Markdown.Header1
	pallete[1] = tmp.General.Tab
	pallete[2] = tmp.General.SelectedTile
	pallete[3] = tmp.General.Error
	pallete[4] = tmp.General.Message
	pallete[5] = tmp.General.Command
	pallete[6] = tmp.Markdown.CodeText
	pallete[7] = tmp.Code.Comment
	pallete[8] = tmp.General.BottomBar

	return pallete, nil
}
