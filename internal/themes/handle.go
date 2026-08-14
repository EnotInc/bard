package themes

import (
	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/screen"
)

func (t *Themes) j() {
	if t.cursor != len(t.list)-1 {
		t.cursor += 1
	}
}

func (t *Themes) k() {
	if t.cursor != 0 {
		t.cursor -= 1
	}
}

func (t *Themes) change() {
	cfg := config.GetConfig()
	name := t.list[t.cursor]

	msg := config.ChangeTheme(name)
	if msg != "" {
		t.SetError(msg)
		return
	}
	cfg.ThemeName = name
	screen.SendCall(calls.PurgeCache)
	t.changeTheme()
}

func (t *Themes) close() {
	screen.SendCall(calls.ClosePopups)
}
