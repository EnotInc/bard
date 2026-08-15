package settings

import (
	"github.com/EnotInc/Bard/internal/enums/cursor"
	"github.com/EnotInc/Bard/internal/enums/keys"
)

func (s *Settings) DrawLineAt(index int) string {
	return s.RenderSettingAt(index)
}

func (s *Settings) Handle(key rune) {
	switch key {
	case 'j':
		s.j()
	case 'k':
		s.k()
	case keys.Enter:
		s.toggelSetting()
	case keys.Esc:
		s.exit()
	}
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

func (s *Settings) PreDraw() {
	if !s.uptade {
		return
	}

	s.uptade = false
	s.updateSettings()
}
