package setting

import (
	"github.com/EnotInc/Bard/internal/enums/cursor"
)

func (s *Settings) DrawLineAt(index int) string {
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
