package tui

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/buffers"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/services"
)

// This function is used to add visual highlight to the selected lines
func (ui *TUI) AddVisual(curMode mode.Mode, l []rune, i int, startOffset, startLine, endOffset, endLine int, Type buffers.BufferType) string {
	var line []rune

	if len(l) == 0 { // if line is empty, returning selected 'new line' symbol
		return string(ui.WithEndLine(string(l)))
	}

	clear := services.ReplaceTabs(l)
	switch curMode {
	case mode.Visual:
		startOffset += services.CursorShiftAt(l, startOffset)
		endOffset += services.CursorShiftAt(l, endOffset)

		if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
			startLine, endLine = endLine, startLine
			startOffset, endOffset = endOffset, startOffset
			if len(l) > 0 {
				endOffset += 1 // too highlight the whole char
			}
		}

		var rendered string
		if Type != buffers.Other {
			rendered = ui.render.Render(l, i, true, true, i == startLine, ui.XScroll, Type)
		} else {
			rendered = string(clear)
		}

		if startLine == i && i == endLine {
			selected := ui.paint(clear[startOffset:endOffset])
			before := services.VisibleSubString(rendered, 0, startOffset-1)
			after := services.VisibleSubString(rendered, endOffset, len(clear))
			line = []rune(before + ascii.Reset.Str() + string(selected) + after)

		} else if startLine < i && i < endLine {
			line = ui.WithEndLine(string(ui.paint(clear)))

		} else if startLine == i {
			selected := ui.paint(clear[startOffset:])
			before := services.VisibleSubString(rendered, 0, startOffset-1)
			line = ui.WithEndLine(before + ascii.Reset.Str() + string(selected))

		} else if endLine == i {
			selected := ui.paint(clear[:endOffset])
			after := services.VisibleSubString(rendered, endOffset, len(clear))
			line = []rune(string(selected) + after)

		} else {
			line = ui.WithEndLine(rendered)
		}

	case mode.Visual_line:
		if startLine > endLine {
			startLine, endLine = endLine, startLine
		}
		if Type != buffers.Other {
			ui.render.Render(l, i, true, true, i == startLine, ui.XScroll, Type)
		}

		theme := config.GetTheme().General
		l := theme.Selection + string(clear) + ascii.Reset.Str()
		line = ui.WithEndLine(l)
	}

	return string(line)
}

// used to add 'new line' symbol to the givven selected line
func (ui *TUI) WithEndLine(l string) []rune {
	theme := config.GetTheme().General
	return []rune(l + theme.Selection + ascii.NewLine.Str() + ascii.Reset.Str())
}

func (ui *TUI) paint(line []rune) []rune {
	var s strings.Builder
	theme := config.GetTheme().General
	s.WriteString(theme.Selection)
	s.WriteString(string(line))
	s.WriteString(ascii.Reset.Str())
	return []rune(s.String())
}
