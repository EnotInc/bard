package themes

import (
	"sort"

	"github.com/EnotInc/Bard/config"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

type ThemeEntry struct {
	name    string
	pallete [9]string
}

type Themes struct {
	list        []ThemeEntry
	SetError    func(msg string)
	changeMode  func(mode mode.Mode)
	changeTheme func()
	cursor      int
	w, h        int
}

func InitThemes(purgeCache func(), changeMode func(mode.Mode), SetError func(msg string)) *Themes {

	var cursor int = 0
	var list []ThemeEntry

	t, err := config.GetThemeList()
	if err != nil {
		list = []ThemeEntry{}
	}

	for n, p := range t {
		entry := ThemeEntry{}
		entry.name = n
		entry.pallete = p
		list = append(list, entry)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].name < list[j].name
	})

	return &Themes{
		changeMode:  changeMode,
		changeTheme: purgeCache,
		SetError:    SetError,
		cursor:      cursor,
		list:        list,
	}
}
