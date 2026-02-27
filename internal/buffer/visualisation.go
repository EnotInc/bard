package buffer

import (
	"Enot/Bard/internal/ascii"
	"Enot/Bard/internal/mode"
	"slices"
)

// This function is used to add visual highlight to the selected lines
func (b *Buffer) AddVisual(curMode mode.Mode, l []rune, i int) string {
	var line []rune

	switch curMode {
	case mode.Visual:
		startOffset := b.Visual.Offset
		startLine := b.Visual.Line

		endOffset := b.Cursor.Offset
		endLine := b.Cursor.Line

		if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
			startLine, endLine = endLine, startLine
			startOffset, endOffset = endOffset, startOffset
		}

		if len(b.Lines[endLine].Data) > 0 {
			endOffset++
		}

		if startLine == i && i == endLine {
			//line = l[:startOffset] + string(startSel) + l[startOffset:endOffset] + string(reset) + l[endOffset:]
			line = slices.Concat(l[:startOffset], []rune(ascii.StartSel), l[startOffset:endOffset], []rune(ascii.Reset), l[endOffset:])
		} else if startLine < i && i < endLine {
			//line = string(startSel) + l + string(reset)
			line = slices.Concat([]rune(ascii.StartSel), l, []rune(ascii.Reset))
		} else if startLine == i {
			//line = l[:startOffset] + string(startSel) + l[startOffset:] + string(reset)
			line = slices.Concat(l[:startOffset], []rune(ascii.StartSel), l[startOffset:], []rune(ascii.Reset))
		} else if endLine == i {
			//line = string(startSel) + l[:endOffset] + string(reset) + l[endOffset:]
			line = slices.Concat([]rune(ascii.StartSel), l[:endOffset], []rune(ascii.Reset), l[endOffset:])
		} else {
			line = l
		}

	case mode.Visual_line:
		startLine := b.Visual.Line
		endLine := b.Cursor.Line

		if startLine > endLine {
			startLine, endLine = endLine, startLine
		}

		line = slices.Concat([]rune(ascii.StartSel), l, []rune(ascii.Reset))
	}

	return string(line)
}
