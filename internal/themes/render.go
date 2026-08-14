package themes

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
)

func (t *Themes) rednerNameAt(index int, offset int) string {
	var name strings.Builder
	theme := config.GetTheme()

	if index == 0 {
		header := " Available themes"
		name.WriteString(theme.General.BottomBar)
		name.WriteString(theme.General.Message)
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
		name.WriteString(theme.General.BottomBar)
		name.WriteString(emtpy)
	} else {
		item := t.list[index]
		n := item.name
		if item.name == config.GetConfig().ThemeName {
			n = fmt.Sprintf("> %s <", n)
		} else {
			n = fmt.Sprintf("  %s", n)
		}

		emtpy := strings.Repeat(" ", max(0, offset-len(n)))

		if index == t.cursor {
			name.WriteString(ascii.UnderLine.Str())
			name.WriteString(ascii.OverLine.Str())
		}

		name.WriteString(item.pallete[8])
		name.WriteString(n)
		name.WriteString(emtpy)
	}

	name.WriteString(ascii.Reset.Str())
	return services.VisibleSubString(name.String(), 0, offset)
}

func (t *Themes) renderPreviewAt(index int, offset int) string {
	var preview strings.Builder
	theme := config.GetTheme()

	if index == 0 {
		preview.WriteString(theme.General.BottomBar)
		preview.WriteString(ascii.UnderLine.Str())
		preview.WriteString(theme.General.Message)
		preview.WriteString(strings.Repeat(" ", max(0, t.w-offset)))
		preview.WriteString(ascii.Reset.Str())
		return services.VisibleSubString(preview.String(), 0, t.w-offset)
	}

	index -= 1 // removeing header offset
	if index >= len(t.list) {
		preview.WriteString(theme.General.BottomBar)
		preview.WriteString(strings.Repeat(" ", max(0, t.w-offset)))
		return services.VisibleSubString(preview.String(), 0, t.w-offset)
	}

	item := t.list[index]

	if index == t.cursor {
		preview.WriteString(ascii.UnderLine.Str())
		preview.WriteString(ascii.OverLine.Str())
	}

	preview.WriteString(item.pallete[8])
	preview.WriteString(strings.Repeat(" ", max(0, t.w-offset-16))) // 16 - amount of pallete symbols
	for i, c := range item.pallete {
		if i == len(item.pallete)-1 {
			continue
		}
		preview.WriteString(c)
		preview.WriteString(ascii.ColorBox.Str())
		preview.WriteString(" ")
	}

	preview.WriteString(ascii.Reset.Str())
	return services.VisibleSubString(preview.String(), 0, t.w-offset)
}
