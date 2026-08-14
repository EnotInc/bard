package themes

import mode "github.com/EnotInc/Bard/internal/enums/mode"

type Themes struct {
	SetError    func(msg string)
	list        []string
	changeMode  func(mode mode.Mode)
	changeTheme func()
	cursor      int
	w, h        int
}

func InitThemes(purgeCache func(), changeMode func(mode.Mode), SetError func(msg string)) *Themes {

	var cursor int = 0
	var list []string = []string{}

	return &Themes{
		changeMode:  changeMode,
		changeTheme: purgeCache,
		SetError:    SetError,
		cursor:      cursor,
		list:        list,
	}
}
