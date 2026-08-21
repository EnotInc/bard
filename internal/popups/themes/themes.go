package themes

import (
	"sort"

	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/theme"
)

type ThemeEntry struct {
	name    string
	pallete [9]string
}

type action int

const (
	_ action = iota
	none
	search
)

type Themes struct {
	list       []ThemeEntry
	searched   []ThemeEntry
	SetError   func(msg string)
	changeMode func(mode mode.Mode)
	onChange   func()
	search     []rune
	action     action
	cursor     int
	w, h       int
}

func InitThemes(purgeCache func(), changeMode func(mode.Mode), SetError func(msg string)) *Themes {

	var cursor int = 0
	var list []ThemeEntry

	t, err := theme.GetThemeList()
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
		changeMode: changeMode,
		onChange:   purgeCache,
		SetError:   SetError,
		cursor:     cursor,
		action:     none,
		list:       list,
	}
}
