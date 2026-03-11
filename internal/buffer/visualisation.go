package buffer

import (
	"slices"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/mode"
)

// This function is used to add visual highlight to the selected lines
func (b *Buffer) AddVisual(curMode mode.Mode, l []rune, i int) string {
	var line []rune

	switch curMode {
	case mode.Visual:
		startOffset := b.Visual.offset
		startLine := b.Visual.line

		endOffset := b.Cursor.offset
		endLine := b.Cursor.line

		if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
			startLine, endLine = endLine, startLine
			startOffset, endOffset = endOffset, startOffset
		}

		if len(b.Lines[endLine].Data) > 0 {
			endOffset++
		}

		if startLine == i && i == endLine {
			line = slices.Concat(l[:startOffset], []rune(ascii.StartSel), l[startOffset:endOffset], []rune(ascii.Reset), l[endOffset:])
		} else if startLine < i && i < endLine {
			line = slices.Concat([]rune(ascii.StartSel), l, []rune(ascii.Reset))
		} else if startLine == i {
			line = slices.Concat(l[:startOffset], []rune(ascii.StartSel), l[startOffset:], []rune(ascii.Reset))
		} else if endLine == i {
			line = slices.Concat([]rune(ascii.StartSel), l[:endOffset], []rune(ascii.Reset), l[endOffset:])
		} else {
			line = l
		}

	case mode.Visual_line:
		startLine := b.Visual.line
		endLine := b.Cursor.line

		if startLine > endLine {
			startLine, endLine = endLine, startLine
		}

		line = slices.Concat([]rune(ascii.StartSel), l, []rune(ascii.Reset))
	}

	return string(line)
}
