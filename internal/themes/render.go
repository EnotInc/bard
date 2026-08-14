package themes

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
)

func (t *Themes) rednerNameAt(index int, offset int) string {
	var name strings.Builder

	if index == 0 {
		header := " Available themes"
		name.WriteString(ascii.Bold.Str())
		name.WriteString(ascii.UnderLine.Str())
		name.WriteString(header)
		name.WriteString(strings.Repeat(" ", max(0, offset-len(header))))
		name.WriteString(ascii.Reset.Str())
		return services.VisibleSubString(name.String(), 0, offset)
	}

	index -= 1 // removeing header offset
	if index >= len(t.list) {
		emtpy := strings.Repeat(" ", int(offset))
		name.WriteString(emtpy)
	} else {
		item := t.list[index]
		_w := max(0, offset-len(item)-2) // -2 - saved space for offset
		emtpy := strings.Repeat(" ", _w)

		if item == config.GetConfig().ThemeName {
			name.WriteString(config.GetTheme().Markdown.Highlight)
		}

		if index == t.cursor {
			name.WriteString(ascii.UnderLine.Str())
			name.WriteString(ascii.OverLine.Str())
		}

		name.WriteString("  ")
		name.WriteString(item)
		name.WriteString(emtpy)
	}

	name.WriteString(ascii.Reset.Str())
	return services.VisibleSubString(name.String(), 0, offset)
}

func (t *Themes) renderPreviewAt(index int, offset int) string {
	var preview strings.Builder

	if index == 0 {
		preview.WriteString(ascii.UnderLine.Str())
		preview.WriteString(strings.Repeat(" ", max(0, t.w-offset)))
		preview.WriteString(ascii.Reset.Str())
		return services.VisibleSubString(preview.String(), 0, t.w-offset)
	}

	index -= 1 // removeing header offset
	if index >= len(t.list) {
		preview.WriteString(strings.Repeat(" ", max(0, t.w-offset)))
		return services.VisibleSubString(preview.String(), 0, t.w-offset)
	}

	pallete, err := config.GetThemePallete(t.list[index])
	if err != nil {
		t.SetError(err.Error())
		return strings.Repeat(" ", t.w-offset)
	}

	if index == t.cursor {
		preview.WriteString(ascii.UnderLine.Str())
		preview.WriteString(ascii.OverLine.Str())
	}

	preview.WriteString(strings.Repeat(" ", max(0, t.w-offset-16))) // 16 - amount of pallete symbols
	for _, c := range pallete {
		preview.WriteString(c)
		preview.WriteString(ascii.ColorBox.Str())
		preview.WriteString(" ")
	}

	return services.VisibleSubString(preview.String(), 0, t.w-offset)
}
