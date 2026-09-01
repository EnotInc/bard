package setting

import (
	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/screen"
)

func (s *Settings) handle(key rune) {
	switch key {
	case 'j':
		s.j()
	case 'k':
		s.k()
	case keys.Enter, keys.Space:
		s.toggle()
	case keys.Esc:
		screen.SendCall(calls.ClosePopups)
	}
}

func (s *Settings) j() {
	if s.cursor < settings_amount {
		s.cursor += 1
	}
}

func (s *Settings) k() {
	if s.cursor > 0 {
		s.cursor -= 1
	}
}

func (s *Settings) toggle() {
	cfg := config.GetConfig()

	switch s.cursor + header_offset {
	case int(KeepTabs):
		cfg.KeepTabs = !cfg.KeepTabs
	case int(RelativeNumbers):
		cfg.RLN = !cfg.RLN
	case int(ShowMDSymbols):
		cfg.ShowMD = !cfg.ShowMD
	case int(ShowTabNames):
		cfg.TabNames = !cfg.TabNames
	case int(EnableRender):
		cfg.Render = !cfg.Render
	case int(ShowDotFiles):
		cfg.ShowDot = !cfg.ShowDot
	case int(ShowBorders):
		cfg.ShowBorder = !cfg.ShowBorder
	case int(ShowIcons):
		cfg.ShowIcons = !cfg.ShowIcons
	case int(ShowEmpty):
		cfg.ShowEmpty = !cfg.ShowEmpty
	}

	s.onChange()
	screen.SendCall(calls.PurgeCache)
}
