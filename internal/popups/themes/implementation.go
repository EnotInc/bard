package themes

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/cursor"
	"github.com/EnotInc/Bard/theme"
)

const sep float32 = 0.5

func (t *Themes) DrawLineAt(index int) string {
	var data strings.Builder

	if index == 0 {
		data.WriteString(t.buildSearchBar())
		data.WriteString(ascii.Reset.Str())
		return data.String()
	}

	theme := theme.GetTheme().General
	cfg := config.GetConfig()
	if (index == t.h-1 && !cfg.ShowBorder) || (index == t.h-3 && cfg.ShowBorder) {
		text := " <k/j> up/down   <enter> select   <esc> exit "
		data.WriteString(theme.BottomBar)
		data.WriteString(theme.Message)
		data.WriteString(ascii.Italic.Str())
		data.WriteString(text)
		data.WriteString(strings.Repeat(" ", max(0, t.w-len(text))))
		data.WriteString(ascii.Reset.Str())
		return data.String()
	}

	offset := float32(t.w) * sep
	name := t.rednerNameAt(index, int(offset))
	preview := t.renderPreviewAt(index, int(offset))

	data.WriteString(name)
	data.WriteString(preview)
	data.WriteString(ascii.Reset.Str())

	return data.String()
}

func (t *Themes) Handle(key rune) {
	switch t.action {
	case search:
		t.handleSearch(key)
	case none:
		t.handleDefault(key)
	}
}

func (t *Themes) GetCursor(withBorder bool) (int, int, cursor.CursorType) {
	var x, y int
	switch t.action {
	case none:
		x = -1
		y = -1
	case search:
		x = searchOffset + len(t.search) + 1
		y = 0

		if withBorder {
			x += 1
			y += 1
		}
	}

	return x, y, cursor.CursorBloc
}

func (t *Themes) SetTitle() string {
	return "Themes"
}

func (t *Themes) Resize(w, h int) {
	t.w = w
	t.h = h
}

func (t *Themes) PreDraw() {
	if len(t.search) == 0 {
		t.searched = t.list
		return
	}

	t.searched = []ThemeEntry{}
	for _, theme := range t.list {
		if strings.Contains(theme.name, string(t.search)) {
			t.searched = append(t.searched, theme)
		}
	}
}
