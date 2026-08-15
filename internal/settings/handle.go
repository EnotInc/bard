package settings

import (
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/screen"
)

func (s *Settings) j() {
	if s.cursor < len(s.bools)-1 {
		s.cursor += 1
	}
}

func (s *Settings) k() {
	if s.cursor > 0 {
		s.cursor -= 1
	}
}

func (s *Settings) exit() {
	s.uptade = true
	screen.SendCall(calls.ClosePopups)
}

func (s *Settings) toggelSetting() {
	s.uptade = true
	switch s.cursor {
	case 0:
		s.bools[EnableRender] = !s.bools[EnableRender]
	case 1:
		s.bools[ShowEmpty] = !s.bools[ShowEmpty]
	case 2:
		s.bools[ShowMDSymbols] = !s.bools[ShowMDSymbols]
	case 3:
		s.bools[ShowTabNames] = !s.bools[ShowTabNames]
	case 4:
		s.bools[ShowIcons] = !s.bools[ShowIcons]
	case 5:
		s.bools[ShowBorders] = !s.bools[ShowBorders]
	case 6:
		s.bools[ShowDotFiles] = !s.bools[ShowDotFiles]
	case 7:
		s.bools[RelativeNumbers] = !s.bools[RelativeNumbers]
	default:
		s.uptade = false
	}

	if s.uptade {
		s.apply()
	}
	screen.SendCall(calls.PurgeCache)
	s.change()
}
