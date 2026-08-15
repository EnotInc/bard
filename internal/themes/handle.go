package themes

import (
	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"
)

func (t *Themes) j() {
	if t.cursor != len(t.searched)-1 {
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
	name := t.searched[t.cursor].name

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

func (t *Themes) handleSearch(key rune) {
	switch key {
	case keys.Esc, keys.Enter:
		t.action = none
		t.cursor = 0
	case keys.Backspace:
		if len(t.search) > 0 {
			t.search = t.search[:len(t.search)-1]
		} else {
			t.action = none
			t.search = []rune{}
		}
	default:
		if services.IsLetterOrNumber(key) {
			t.search = append(t.search, key)
		}
	}
}

func (t *Themes) handleDefault(key rune) {
	switch key {
	case ':':
		t.changeMode(mode.Command)
		t.close()
	case 'j': // down
		t.j()
	case 'k': // up
		t.k()
	case '/':
		t.action = search
		t.cursor = -1
	case keys.Enter:
		t.change()
	case keys.Esc:
		t.close()
	}
}
