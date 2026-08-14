package themes

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
)

func (t *Themes) rednerNameAt(index int, offset int) string {
	var name strings.Builder
	if index >= len(t.list) {
		emtpy := strings.Repeat(" ", int(offset))
		name.WriteString(emtpy)
	} else {
		item := t.list[index]
		_w := max(0, offset-len(item)-2) // -2 - saved space for cursor
		emtpy := strings.Repeat(" ", _w)

		if item == config.GetConfig().ThemeName {
			name.WriteString(config.GetTheme().Markdown.Highlight)
		}

		// cursor
		if index == t.cursor {
			name.WriteString("> ")
		} else {
			name.WriteString("  ")
		}

		name.WriteString(item)
		name.WriteString(emtpy)
	}

	name.WriteString(ascii.Reset.Str())
	return services.VisibleSubString(name.String(), 0, offset)
}

func (t *Themes) renderPreviewAt(index int, offset int) string {
	var preview strings.Builder
	var text string = ""
	theme := config.GetTheme()
	switch index {
	case 0:
		preview.WriteString(theme.Markdown.Header1)
		text = "Header 1"
	case 1:
		preview.WriteString(theme.Markdown.Image)
		text = "Image"
	case 2:
		preview.WriteString(theme.Markdown.CodeText)
		text = "CodeLine"
	case 3:
		preview.WriteString(theme.Markdown.Quote)
		text = "Quote"
	case 4:
		preview.WriteString(theme.Markdown.NumberList)
		text = "NumberList"
	case 5:
		preview.WriteString(theme.General.SelectedTile)
		text = "SelectedTile"
	case 6:
		preview.WriteString(theme.General.CurrentLine)
		text = "CurrentLine"
	case 7:
		preview.WriteString(theme.General.EmptyLine)
		text = "EmptyLine"
	case 8:
		preview.WriteString(theme.General.Command)
		text = "Command"
	case 9:
		preview.WriteString(theme.General.Error)
		text = "Error"
	case 10:
		preview.WriteString(theme.Code.Symbol)
		text = "Symbol"
	case 11:
		preview.WriteString(theme.Code.Keyword)
		text = "Keyword"
	case 12:
		preview.WriteString(theme.Code.Bracket)
		text = "Brackets"
	case 13:
		preview.WriteString(theme.Code.String)
		text = "String"
	case 14:
		preview.WriteString(theme.Code.Number)
		text = "Number"
	}
	if len(text) != 0 {
		preview.WriteByte(' ')
		preview.WriteString(ascii.ColorBox.Str())
		preview.WriteByte(' ')
	} else {
		preview.WriteString("   ")
	}
	preview.WriteString(ascii.ResetFg.Str())
	preview.WriteString(text)
	preview.WriteString(strings.Repeat(" ", max(0, t.w-offset-len(text)-3))) // -3 - offset for symbols

	return services.VisibleSubString(preview.String(), 0, t.w-offset)
}
