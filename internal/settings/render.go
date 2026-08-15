package settings

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
)

func (s *Settings) RenderSettingAt(index int) string {
	var data strings.Builder
	var text string
	var enable bool

	theme := config.GetTheme()
	cfg := config.GetConfig()

	data.WriteString(theme.General.BottomBar)
	data.WriteString(" ")

	switch index {
	case 0:
		text = "EnableRender"
		enable = s.bools[EnableRender]
	case 1:
		text = "ShowEmpty"
		enable = s.bools[ShowEmpty]
	case 2:
		text = "ShowMDSymbols"
		enable = s.bools[ShowMDSymbols]
	case 3:
		text = "ShowTabNames"
		enable = s.bools[ShowTabNames]
	case 4:
		text = "ShowIcons"
		enable = s.bools[ShowIcons]
	case 5:
		text = "ShowBorders"
		enable = s.bools[ShowBorders]
	case 6:
		text = "ShowDotFiles"
		enable = s.bools[ShowDotFiles]
	case 7:
		text = "RelativeNumbers"
		enable = s.bools[RelativeNumbers]
	default:
		text = ""
	}

	if text != "" {
		switch enable {
		case true:
			data.WriteString(theme.General.Message)
		case false:
			data.WriteString(theme.General.Error)
		}

		icon := services.SwitchIcon(enable, cfg.ShowIcons)
		data.WriteString(icon)
		data.WriteString(ascii.ResetFg.Str())
	} else {
		data.WriteString("  ")
	}

	if s.cursor == index {
		data.WriteString(ascii.OverLine.Str())
		data.WriteString(ascii.UnderLine.Str())
	}

	data.WriteString(text)
	data.WriteString(ascii.Reset.Str())
	data.WriteString(theme.General.BottomBar)
	data.WriteString(strings.Repeat(" ", max(0, s.w-len(text)-2))) // -2 - switch icon offset

	return data.String()
}
