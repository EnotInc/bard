package themes

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/cursor"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

const sep float32 = 0.5

func (t *Themes) DrawLineAt(index int) string {
	var data strings.Builder

	theme := config.GetTheme().General

	if index == t.h-1 {
		text := " k/j - up/down   <enter> - select   <esc> - exit"
		data.WriteString(theme.BottomBar)
		data.WriteString(theme.Message)
		data.WriteString(text)
		data.WriteString(strings.Repeat(" ", max(0, t.w-len(text))))
		data.WriteString(ascii.Reset.Str())
		return data.String()
	}

	offset := float32(t.w) * sep
	name := t.rednerNameAt(index, int(offset))
	preview := t.renderPreviewAt(index, int(offset))

	data.WriteString(theme.BottomBar)
	data.WriteString(name)
	data.WriteString(theme.BottomBar)
	data.WriteString(preview)
	data.WriteString(ascii.Reset.Str())

	return data.String()
}

func (t *Themes) Handle(key rune) {
	switch key {
	case ':':
		t.changeMode(mode.Command)
		t.close()
	case 'j': // down
		t.j()
	case 'k': // up
		t.k()
	case keys.Enter:
		t.change()
	case keys.Esc:
		t.close()
	}
}

func (t *Themes) GetCursor(withBorder bool) (int, int, cursor.CursorType) {
	return 1, 1, cursor.CursorBloc
}

func (t *Themes) SetTitle() string {
	return "Themes"
}

func (t *Themes) Resize(w, h int) {
	t.w = w
	t.h = h
}

func (t *Themes) PreDraw() {
	var err error
	t.list, err = config.GetThemeList()
	if err != nil {
		t.SetError(err.Error())
	}
}
