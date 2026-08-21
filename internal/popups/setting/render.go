package setting

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
	"github.com/EnotInc/Bard/theme"
)

func (s *Settings) render(index int) string {
	var data strings.Builder

	var icon string
	var text string
	var enable bool

	cfg := config.GetConfig()
	theme := theme.GetTheme().General

	switch index {
	case 0:
		text = "Settings"
	case int(RelativeNumbers):
		text = "Relative Line Number"
		enable = cfg.RLN
	case int(ShowMDSymbols):
		text = "Show MarkDonw Symbols"
		enable = cfg.ShowMD
	case int(ShowTabNames):
		text = "Show Tab Names"
		enable = cfg.TabNames
	case int(EnableRender):
		text = "Enable Render"
		enable = cfg.Render
	case int(ShowDotFiles):
		text = "Show Dot Files"
		enable = cfg.ShowDot
	case int(ShowBorders):
		text = "Show Borders"
		enable = cfg.ShowBorder
	case int(ShowIcons):
		text = "Show Icons"
		enable = cfg.ShowIcons
	case int(ShowEmpty):
		text = "Show Empty Lines"
		enable = cfg.ShowEmpty
	default:
		text = ""
	}

	const iconOffset int = 3
	var offset int = 0
	if cfg.ShowBorder {
		offset = 1
	}

	if text == "" || index == 0 {
		icon = "   "
	} else {
		icon = " " + services.SwitchIcon(enable, cfg.ShowIcons)
	}

	if index == 0 {
		data.WriteString(ascii.UnderLine.Str())
		data.WriteString(ascii.Bold.Str())
	}

	data.WriteString(theme.BottomBar)
	data.WriteString(icon)
	data.WriteString(ascii.ResetFg.Str())

	if index == s.cursor+header_offset {
		data.WriteString(ascii.UnderLine.Str())
		data.WriteString(ascii.OverLine.Str())
	}

	data.WriteString(text)
	data.WriteString(strings.Repeat(" ", max(0, s.w-len(text)-iconOffset-offset)))
	data.WriteString(ascii.Reset.Str())

	return data.String()
}
