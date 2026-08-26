package themes

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
	"github.com/EnotInc/Bard/theme"
)

func (t *Themes) renderAt(index int) string {
	var data strings.Builder
	theme := theme.GetTheme()
	cfg := config.GetConfig()

	index -= 1 // removing search bard offset

	if index >= len(t.searched) { // on empty lines
		empty := strings.Repeat(" ", t.w)
		data.WriteString(theme.General.BottomBar)
		data.WriteString(empty)
	} else {
		const palleteOffset int = 16
		const spacing int = 6
		item := t.searched[index]
		name := item.name
		clear := strings.TrimSuffix(name, ".json")

		if item.name == cfg.ThemeName {
			name = fmt.Sprintf("> %s <", clear)
		} else {
			name = fmt.Sprintf("  %s", clear)
		}

		var borderOffset int = 0
		if cfg.ShowBorder {
			borderOffset = 2
		}

		amount := services.CountClear(name, 0, len(name))
		empty := strings.Repeat(" ", max(0, t.w-amount-palleteOffset-borderOffset-spacing))

		if index == t.cursor {
			data.WriteString(ascii.UnderLine.Str())
			data.WriteString(ascii.OverLine.Str())
		}

		border := fmt.Sprintf("%s%s", theme.General.BottomBar, strings.Repeat(" ", spacing/2))

		data.WriteString(border)
		data.WriteString(item.pallete[8])
		data.WriteString(name)
		data.WriteString(empty)

		for i, c := range item.pallete {
			if i == len(item.pallete)-1 {
				continue
			}
			data.WriteString(c)
			data.WriteString(ascii.ColorBox.Str())
		}
		data.WriteString(border)
	}

	data.WriteString(ascii.Reset.Str())
	return data.String()
}
