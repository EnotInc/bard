package buffer

import "slices"

type copied struct {
	data    []rune
	isEnd   bool
	isStart bool
}

func (b *Buffer) SwapTail() {
	b.Visual.Offset, b.Cursor.Offset = b.Cursor.Offset, b.Visual.Offset
	b.Visual.Line, b.Cursor.Line = b.Cursor.Line, b.Visual.Line

}

// Create a new line for the buffer.copied list
func (b *Buffer) CopyLine(l *Line, startOffset int, endOffset int) *copied {
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

func (b *Buffer) CopySelected(isDelete bool, isVisualLine bool) {
	b.copies = []*copied{}

	startOffset := b.Visual.Offset
	startLine := b.Visual.Line
	endOffset := b.Cursor.Offset
	endLine := b.Cursor.Line

	if startLine > endLine || (startLine == endLine && startOffset > endOffset) {
		startLine, endLine = endLine, startLine
		startOffset, endOffset = endOffset, startOffset
	}

	var tempLine []rune

	lineCount := 0
	lineSelected := endLine - startLine
	for i := startLine; i <= endLine; {
		curOffsetStart := 0
		if lineCount == 0 {
			curOffsetStart = startOffset
		}

		curOffsetEnd := max(len(b.Lines[i].Data)-1, 0)
		if i == endLine && !isVisualLine {
			curOffsetEnd = endOffset
		}
		if len(b.Lines[i].Data) > 0 {
			curOffsetEnd++
		}

		line := b.CopyLine(b.Lines[i], curOffsetStart, curOffsetEnd)
		b.copies = append(b.copies, line)

		if isDelete {
			if lineCount == 0 {
				tempLine = b.Lines[i].Data[:curOffsetStart]
			}
			if lineCount == lineSelected {
				tempLine = append(tempLine, b.Lines[i].Data[curOffsetEnd:]...)
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
	b.Cursor.Line = startLine
	if b.Cursor.Line > len(b.Lines)-1 {
		b.Cursor.Line = max(len(b.Lines)-1, 0)
	}

	b.Cursor.Offset = startOffset
	b.fixOffset()
}

func (b *Buffer) Paste(shift int) {
	initialOffset := b.Cursor.Offset + shift
	if len(b.Lines[b.Cursor.Line].Data) == 0 || b.Cursor.Offset < 0 {
		initialOffset = 0
	}

	DataFirst := append([]rune(nil), b.Lines[b.Cursor.Line].Data[:initialOffset]...)
	DataSecond := append([]rune(nil), b.Lines[b.Cursor.Line].Data[initialOffset:]...)

	isFirstStart := b.copies[0].isStart
	isLastEnd := b.copies[len(b.copies)-1].isEnd

	lineIndex := b.Cursor.Line
	for i, line := range b.copies {
		Data := append([]rune{}, line.data...)

		lineIndex = b.Cursor.Line + i // Moving index while walking on copied lines
		if lineIndex >= len(b.Lines) {
			lineIndex = len(b.Lines) - 1
		}
		curLine := b.Lines[lineIndex]

		switch i {
		case 0: // Working with 1st line
			if len(b.copies) == 1 {
				if !isFirstStart && !isLastEnd {
					curLine.Data = slices.Concat(DataFirst, Data, DataSecond)
				} else {
					b.InsertLineWithData(lineIndex+shift, Data)
				}
			} else {
				if isFirstStart && isLastEnd {
					b.InsertLineWithData(lineIndex+shift, Data)
				} else {
					curLine.Data = slices.Concat(DataFirst, Data)
				}
			}
		case len(b.copies) - 1: // Working with last line
			if isLastEnd && isFirstStart {
				b.InsertLineWithData(lineIndex+shift, Data)
			} else {
				savedData := slices.Concat(Data, DataSecond)
				b.InsertLineWithData(lineIndex+shift, savedData)
			}
		default: //Just insert a new line
			b.InsertLineWithData(lineIndex+shift, Data)
		}
	}

	b.fixOffset()
}
