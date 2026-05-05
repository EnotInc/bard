package tui

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
)

// About AddVisual()
// This function is used to add visual highlight to the selected lines
func (ui *TUI) AddVisual(curMode enums.Mode, l []rune, i int, startOffset, startLine, endOffset, endLine int, lastLineLen int) string {
	var line []rune

	if len(l) == 0 { // if line is empty, returning selected 'new line' symbol
		return ui.theme.Selection + ascii.NewLine.Str() + ascii.Reset.Str()
	}

	switch curMode {
	case enums.Visual:
		if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
			startLine, endLine = endLine, startLine
			startOffset, endOffset = endOffset, startOffset
			if len(l) > 0 {
				endOffset += 1 // too highlight the whole char
			}
		} else if lastLineLen > 0 {
			endOffset += 1 // too highlight the whole char
		}

		rendered, _ := ui.render.Render(l, i, true, true, i == startLine)
		if startLine == i && i == endLine {
			selected := ui.paint(l[startOffset:endOffset])
			before := VisibleSubString(rendered, 0, startOffset-1)
			after := VisibleSubString(rendered, endOffset, len(l))
			line = []rune(before + ascii.Reset.Str() + string(selected) + after)

		} else if startLine < i && i < endLine {
			line = ui.WithEndLine(string(ui.paint(l)))

		} else if startLine == i {
			selected := ui.paint(l[startOffset:])
			before := VisibleSubString(rendered, 0, startOffset-1)
			line = ui.WithEndLine(before + ascii.Reset.Str() + string(selected))

		} else if endLine == i {
			selected := ui.paint(l[:endOffset])
			after := VisibleSubString(rendered, endOffset, len(l))
			line = []rune(string(selected) + after)

		} else {
			line = ui.WithEndLine(rendered)
		}

	case enums.Visual_line:
		if startLine > endLine {
			startLine, endLine = endLine, startLine
		}

		l := ui.theme.Selection + string(l) + ascii.Reset.Str()
		line = ui.WithEndLine(l)
	}

	return string(line)
}

// About WithEndLine()
// used to add 'new line' symbol to the givven selected line
func (ui *TUI) WithEndLine(l string) []rune {
	return []rune(l + ui.theme.Selection + ascii.NewLine.Str() + ascii.Reset.Str())
}

// About paint()
// used to colorise every single char in line
// is just inserts selected ascii.StarSel [Color] before the char
func (ui *TUI) paint(line []rune) []rune {
	var s strings.Builder
	for _, ch := range line {
		fmt.Fprintf(&s, "%s%c", ui.theme.Selection, ch)
	}
	s.WriteString(ascii.Reset.Str())
	return []rune(s.String())
}
