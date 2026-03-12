package tui

import (
	"fmt"
	"slices"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/mode"
)

// This function is used to add visual highlight to the selected lines
// TODO: move to TUI
func AddVisual(curMode mode.Mode, l []rune, i int, startOffset, startLine, endOffset, endLine int, lastLineLen int) string {
	var line []rune

	switch curMode {
	case mode.Visual:
		if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
			startLine, endLine = endLine, startLine
			startOffset, endOffset = endOffset, startOffset
		}

		if lastLineLen > 0 {
			endOffset += 1 // too highlight the whole char
		}

		if startLine == i && i == endLine {
			selected := paint(l[startOffset:endOffset])
			line = slices.Concat(l[:startOffset], selected, l[endOffset:])
		} else if startLine < i && i < endLine {
			line = paint(l)
		} else if startLine == i {
			selected := paint(l[startOffset:])
			line = slices.Concat(l[:startOffset], selected)
		} else if endLine == i {
			selected := paint(l[:endOffset])
			line = slices.Concat(selected, l[endOffset:])
		} else {
			line = l
		}

	case mode.Visual_line:
		if startLine > endLine {
			startLine, endLine = endLine, startLine
		}

		line = slices.Concat([]rune(ascii.StartSel), l, []rune(ascii.Reset))
	}

	return string(line)
}

func paint(line []rune) []rune {
	var s = ""
	for _, x := range line {
		s += fmt.Sprintf("%s%c", ascii.StartSel, x)
	}
	s += ascii.Reset.Str()
	return []rune(s)
}
