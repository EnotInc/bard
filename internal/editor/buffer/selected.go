package buffer

import (
	"slices"
	"strings"

	"github.com/EnotInc/Bard/internal/enums"
	cases "github.com/EnotInc/Bard/internal/enums/cases"
)

type copied struct {
	data    []rune
	isEnd   bool
	isStart bool
}

func (b *Buffer) StartVisual() {
	b.Visual.line = b.Cursor.line
	b.Visual.offset = b.Cursor.offset
}

func (b *Buffer) StartVisualLine() {
	b.Visual.line = b.Cursor.line
	b.Visual.offset = 0
}

func (b *Buffer) SwapTail() {
	b.Visual.offset, b.Cursor.offset = b.Cursor.offset, b.Visual.offset
	b.Visual.line, b.Cursor.line = b.Cursor.line, b.Visual.line

}

// Create a new line for the buffer.copied list
func (b *Buffer) CopyLine(l *Line, startOffset int, endOffset int) *copied {
	if b.IsReadOnly {
		return nil
	}

	_isStart := startOffset == 0
	_isEnd := endOffset == len(l.Data)

	if startOffset < 0 {
		startOffset = 0
	}
	if endOffset > len(l.Data)-1 {
		endOffset = len(l.Data)
	}
	if endOffset < 0 {
		endOffset = 0
	}
	if startOffset > endOffset {
		startOffset, endOffset = endOffset, startOffset
	}
	var newData []rune
	newData = l.Data[startOffset:endOffset]
	if len(newData) == 0 {
		newData = []rune("")
	}
	Data := append([]rune{}, newData...)
	return &copied{data: Data, isEnd: _isEnd, isStart: _isStart}
}

// copies selected area into [copied], betven [Cursor] and [Visual] points
func (b *Buffer) CopySelected(isDelete bool, isVisualLine bool) {
	if b.IsReadOnly {
		return
	}
	b.Copies = []*copied{}

	startOffset := b.Visual.offset
	startLine := b.Visual.line
	endOffset := b.Cursor.offset
	endLine := b.Cursor.line

	if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
		startLine, endLine = endLine, startLine
		startOffset, endOffset = endOffset, startOffset
	}

	var tempLine []rune

	lineCount := 0
	lineSelected := endLine - startLine
	for i := startLine; i <= endLine; {
		curOfsetStart := 0
		curOfsetEnd := max(len(b.Lines[i].Data)-1, 0)

		if isVisualLine {
			curOfsetStart = 0
			// not changing curOfsetEnd coz it was set to the max of 0 and len of the current line
		} else {
			if lineCount == 0 {
				curOfsetStart = startOffset
			}
			if i == endLine {
				curOfsetEnd = endOffset
			}
		}

		if len(b.Lines[i].Data) > 0 {
			curOfsetEnd++
		}

		line := b.CopyLine(b.Lines[i], curOfsetStart, curOfsetEnd)
		b.Copies = append(b.Copies, line)

		if isDelete {
			if lineCount == 0 {
				tempLine = b.Lines[i].Data[:curOfsetStart]
			}
			if lineCount == lineSelected {
				tempLine = append(tempLine, b.Lines[i].Data[curOfsetEnd:]...)
			}
			b.RemoveLineAt(i)
			endLine--
		} else {
			i++
		}

		lineCount++
	}

	if isDelete && len(tempLine) > 0 {
		b.InsertLineWithData(startLine, tempLine)
	}
	b.Cursor.line = startLine
	if b.Cursor.line > len(b.Lines)-1 {
		b.Cursor.line = max(len(b.Lines)-1, 0)
	}

	if startOffset > len(b.Lines[b.Cursor.line].Data)-1 {
		startOffset = len(b.Lines[b.Cursor.line].Data) - 1
	}
	b.Cursor.offset = min(startOffset, len(b.Lines[b.Cursor.line].Data)-1)
}

func (b *Buffer) SaveCopied() {
	if len(b.Copies) > 1 {
		b.SaveChanges(
			Insert,
			b.Cursor.Line(),
			len(b.Copies)+b.Cursor.Line()-1,
			enums.Without)

	} else {
		operation := Change
		if b.Copies[0].isEnd && b.Copies[0].isStart {
			operation = Insert
		}
		b.SaveChanges(
			operation,
			b.Cursor.Line(),
			b.Cursor.Line(),
			enums.Without)
	}
}

func (b *Buffer) Paste(shift int) {
	if b.IsReadOnly || len(b.Copies) == 0 {
		return
	}

	initialOffset := b.Cursor.offset + shift
	if len(b.Lines[b.Cursor.line].Data) == 0 || b.Cursor.offset < 0 {
		initialOffset = 0
	}
	if initialOffset >= len(b.Lines[b.Cursor.line].Data) {
		initialOffset = len(b.Lines[b.Cursor.line].Data)
	}

	DataFirst := append([]rune(nil), b.Lines[b.Cursor.line].Data[:initialOffset]...)
	DataSecond := append([]rune(nil), b.Lines[b.Cursor.line].Data[initialOffset:]...)

	isFirstStart := b.Copies[0].isStart
	isLastEnd := b.Copies[len(b.Copies)-1].isEnd

	lineIndex := b.Cursor.line
	for i, line := range b.Copies {
		Data := append([]rune{}, line.data...)

		lineIndex = b.Cursor.line + i // Moving index while walking on copied lines
		if lineIndex >= len(b.Lines) {
			lineIndex = len(b.Lines) - 1
		}
		curLine := b.Lines[lineIndex]

		if isFirstStart && i == 0 || isLastEnd && i == len(b.Copies) {
			b.InsertLineWithData(lineIndex, Data)
			continue
		}

		if !line.isEnd && !line.isStart {
			curLine.Data = slices.Concat(DataFirst, Data, DataSecond)
		} else if line.isEnd && line.isStart {
			b.InsertLineWithData(lineIndex, Data)
		} else if line.isStart {
			curLine.Data = slices.Concat(Data, DataSecond)
		} else if line.isEnd {
			savedData := slices.Concat(DataFirst, Data)
			b.InsertLineWithData(lineIndex, savedData)
		}
	}

	b.FixOffset()
}

func (b *Buffer) ChangeLetterCaseTo(Case cases.Case, isVisualLine bool) {
	if b.IsReadOnly {
		return
	}

	startOffset := b.Visual.offset
	startLine := b.Visual.line
	endOffset := b.Cursor.offset
	endLine := b.Cursor.line

	if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
		startLine, endLine = endLine, startLine
		startOffset, endOffset = endOffset, startOffset
	}

	lineCount := 0
	for i := startLine; i <= endLine; {
		curOfsetStart := 0
		curOfsetEnd := max(len(b.Lines[i].Data)-1, 0)

		if isVisualLine {
			curOfsetStart = 0
			// not changing curOfsetEnd coz it was set to the max of 0 and len of the current line
		} else {
			if lineCount == 0 {
				curOfsetStart = startOffset
			}
			if i == endLine {
				curOfsetEnd = endOffset
			}
		}

		if len(b.Lines[i].Data) > 0 {
			curOfsetEnd++
		}

		var changed_line string
		switch Case {
		case cases.Lower:
			changed_line = strings.ToLower(string(b.Lines[i].Data[curOfsetStart:curOfsetEnd]))
		case cases.Upper:
			changed_line = strings.ToUpper(string(b.Lines[i].Data[curOfsetStart:curOfsetEnd]))
		default:
			changed_line = ""
		}

		b.Lines[i].Data = slices.Concat(b.Lines[i].Data[:curOfsetStart], []rune(changed_line), b.Lines[i].Data[curOfsetEnd:])

		i++
		lineCount++
	}
	b.Cursor.line = startLine
	if b.Cursor.line > len(b.Lines)-1 {
		b.Cursor.line = max(len(b.Lines)-1, 0)
	}

	b.Cursor.offset = startOffset
}
