package setting

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/cursor"
	"github.com/EnotInc/Bard/theme"
)

func (s *Settings) DrawLineAt(index int) string {

	theme := theme.GetTheme().General
	cfg := config.GetConfig()
	if (index == s.h-1 && !cfg.ShowBorder) || (index == s.h-3 && cfg.ShowBorder) {
		var data strings.Builder
		text := " <k/j> up/down   <enter> toggle   <esc> exit "
		data.WriteString(theme.BottomBar)
		data.WriteString(theme.Message)
		data.WriteString(ascii.Italic.Str())
		data.WriteString(text)
		data.WriteString(strings.Repeat(" ", max(0, s.w-len(text))))
		data.WriteString(ascii.Reset.Str())
		return data.String()
	}

	return s.render(index)
}

func (s *Settings) Handle(key rune) {
	s.handle(key)
}

func (s *Settings) GetCursor(withBorder bool) (int, int, cursor.CursorType) {
	return -1, -1, cursor.CursorBloc
}

func (s *Settings) SetTitle() string {
	return ""
}

func (s *Settings) Resize(w, h int) {
	s.w = w
	s.h = h
}

func (s *Settings) PreDraw() {}
