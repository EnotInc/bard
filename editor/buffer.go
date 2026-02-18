package editor

import (
	"slices"
)

type line struct {
	data []rune
}

type copied struct {
	data    []rune
	isEnd   bool
	isStart bool
}

type cursor struct {
	line   int
	offset int

	//keepOffset is uset to keep the ccursor in one place
	// on the X-axis when moving betweeen lines
	keepOffset int
}

type Buffer struct {
	copies []*copied
	lines  []*line
	cursor *cursor
	visual *cursor
}

func InitBuffer() *Buffer {
	c := &cursor{line: 0, offset: 0}
	v := &cursor{line: 0, offset: 0}
	b := &Buffer{
		cursor: c,
		visual: v,
	}
	b.lines = append(b.lines, &line{data: []rune("")})
	return b
}

func (b *Buffer) fixOffset() {
	if b.cursor.offset != b.cursor.keepOffset {
		b.cursor.offset = b.cursor.keepOffset
	}
	if b.cursor.offset > len(b.lines[b.cursor.line].data)-1 {
		b.cursor.offset = len(b.lines[b.cursor.line].data) - 1
	}
	if b.cursor.offset < 0 {
		b.cursor.offset = 0
	}
}

// Move cursor left
func (b *Buffer) H(amount int) {
	if b.cursor.offset-amount > 0 {
		b.cursor.offset -= amount
	} else {
		b.cursor.offset = 0
	}
	b.cursor.keepOffset = b.cursor.offset
}

// Move cursor down
func (b *Buffer) J(amount int) {
	if b.cursor.line+amount < len(b.lines)-1 {
		b.cursor.line += amount
	} else {
		b.cursor.line = len(b.lines) - 1
	}
	b.fixOffset()
}

// Move cursor up
func (b *Buffer) K(amount int) {
	if b.cursor.line-amount > 0 {
		b.cursor.line -= amount
	} else {
		b.cursor.line = 0
	}
	b.fixOffset()
}

// Move cursor right
func (b *Buffer) L(amount int) {
	if b.cursor.offset+amount < len(b.lines[b.cursor.line].data)-1 {
		b.cursor.offset += amount
	} else {
		b.cursor.offset = len(b.lines[b.cursor.line].data) - 1
		b.fixOffset()
	}
	b.cursor.keepOffset = b.cursor.offset
}

func (b *Buffer) InsertKey(key rune) {
	curLine := b.lines[b.cursor.line]
	curLine.data = append(curLine.data[:b.cursor.offset], append([]rune{key}, curLine.data[b.cursor.offset:]...)...)
	b.cursor.offset += 1
}

// Called when the user presses [backspace] and just removes the character in front of it
func (b *Buffer) RemoveKey(keyShift int) {
	if b.cursor.offset > 0 {
		curLine := b.lines[b.cursor.line]
		index := keyShift + b.cursor.offset
		curLine.data = slices.Delete(curLine.data, index-1, index)
		b.cursor.offset -= 1
	} else {
		b.DelAndMoveLine()
	}
}

// Called when the user presses [x] or [s] in normal mode. It deletes the character and copies it to the buffer
func (b *Buffer) Delkey() {
	if len(b.lines[b.cursor.line].data) > 0 {
		curLine := b.lines[b.cursor.line]
		index := b.cursor.offset
		ch := curLine.data[index]
		b.copies = append([]*copied{}, &copied{data: []rune{ch}, isStart: false, isEnd: false})
		curLine.data = slices.Delete(curLine.data, index, index+1)
	}
}

func (b *Buffer) InsertEmptyLine(lineShift int) {
	index := b.cursor.line + lineShift
	newLine := make([]*line, 0)
	newLine = append(newLine, &line{data: []rune("")})
	b.lines = append(b.lines[:index], append(newLine, b.lines[index:]...)...)
}

func (b *Buffer) InsertLineWithData(index int, data []rune) {
	newLine := make([]*line, 0)
	newLine = append(newLine, &line{data: data})
	b.lines = append(b.lines[:index], append(newLine, b.lines[index:]...)...)
}

// Called when the user presses [enter] in the middle of a line. This function shifts data from the right to the new line
func (b *Buffer) InsertLine() {
	index := b.cursor.line + 1
	shiftData := b.lines[b.cursor.line].data[b.cursor.offset:]
	b.lines[b.cursor.line].data = b.lines[b.cursor.line].data[:b.cursor.offset]

	newLine := line{}
	b.lines = append(b.lines[:index], append([]*line{&newLine}, b.lines[index:]...)...)
	b.cursor.line += 1
	b.cursor.offset = 0

	b.lines[b.cursor.line].data = append(b.lines[b.cursor.line].data, shiftData...)
}

// Called when the user deletes the 0th character in a line. The line is deleted and data is moved to the line above
func (b *Buffer) DelAndMoveLine() {
	if b.cursor.line > 0 {
		shiftData := b.lines[b.cursor.line].data[b.cursor.offset:]
		b.RemoveLine()
		b.cursor.line -= 1
		b.cursor.offset = len(b.lines[b.cursor.line].data)
		b.lines[b.cursor.line].data = append(b.lines[b.cursor.line].data, shiftData...)
	}
}

func (b *Buffer) DelAndMoveLineAt(startLine int, endLine int, endOffset int) {
	shiftData := b.lines[endLine].data[endOffset:]
	b.RemoveLineAt(endLine)
	b.cursor.offset = len(b.lines[startLine].data)
	b.lines[startLine].data = append(b.lines[startLine].data, shiftData...)
}

// Delete the whole line
func (b *Buffer) RemoveLine() {
	if len(b.lines) == 1 {
		b.ClearLine()
		return
	}
	b.lines = slices.Delete(b.lines, b.cursor.line, b.cursor.line+1)
	if b.cursor.line >= len(b.lines) {
		b.cursor.line = len(b.lines)
	}
}

// Delete the whole line at index
func (b *Buffer) RemoveLineAt(lineIndex int) {
	if len(b.lines) == 1 {
		b.ClearLine()
		return
	}
	b.lines = slices.Delete(b.lines, lineIndex, lineIndex+1)
	if b.cursor.line >= len(b.lines) {
		b.K(1)
	}
}

// Set line.data = ""
func (b *Buffer) ClearLine() {
	b.lines[b.cursor.line].data = []rune{}
}

// Move buffer.cursor to the first non-space character in the line
func (b *Buffer) moveToFirst() {
	for i := range len(b.lines[b.cursor.line].data) {
		if b.lines[b.cursor.line].data[i] != ' ' {
			b.cursor.offset = i
			break
		}
	}
}

// Create a new line for the buffer.copied list
func (b *Buffer) copyLine(l *line, startOffset int, endOffset int) *copied {
	_isStart := startOffset == 0
	_isEnd := endOffset == len(l.data)

	if startOffset < 0 {
		startOffset = 0
	}
	if endOffset > len(l.data)-1 {
		endOffset = len(l.data)
	}
	if endOffset < 0 {
		endOffset = 0
	}
	if startOffset > endOffset {
		startOffset, endOffset = endOffset, startOffset
	}
	var newData []rune
	newData = l.data[startOffset:endOffset]
	if len(newData) == 0 {
		newData = []rune("")
	}
	data := append([]rune{}, newData...)
	return &copied{data: data, isEnd: _isEnd, isStart: _isStart}
}

func (b *Buffer) copySelected(isDelete bool, isVisualLine bool) {
	b.copies = []*copied{}

	startOffset := b.visual.offset
	startLine := b.visual.line
	endOffset := b.cursor.offset
	endLine := b.cursor.line

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

		curOffsetEnd := max(len(b.lines[i].data)-1, 0)
		if i == endLine && !isVisualLine {
			curOffsetEnd = endOffset
		}
		if len(b.lines[i].data) > 0 {
			curOffsetEnd++
		}

		line := b.copyLine(b.lines[i], curOffsetStart, curOffsetEnd)
		b.copies = append(b.copies, line)

		if isDelete {
			if lineCount == 0 {
				tempLine = b.lines[i].data[:curOffsetStart]
			}
			if lineCount == lineSelected {
				tempLine = append(tempLine, b.lines[i].data[curOffsetEnd:]...)
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
	b.cursor.line = startLine
	if b.cursor.line > len(b.lines)-1 {
		b.cursor.line = max(len(b.lines)-1, 0)
	}

	b.cursor.offset = startOffset
	b.fixOffset()
}

func (b *Buffer) paste(shift int) {
	initialOffset := b.cursor.offset + shift
	if len(b.lines[b.cursor.line].data) == 0 || b.cursor.offset < 0 {
		initialOffset = 0
	}

	dataFirst := append([]rune(nil), b.lines[b.cursor.line].data[:initialOffset]...)
	dataSecond := append([]rune(nil), b.lines[b.cursor.line].data[initialOffset:]...)

	isFirstStart := b.copies[0].isStart
	isLastEnd := b.copies[len(b.copies)-1].isEnd

	lineIndex := b.cursor.line
	for i, line := range b.copies {
		data := append([]rune{}, line.data...)

		lineIndex = b.cursor.line + i // Moving index while walking on copied lines
		if lineIndex >= len(b.lines) {
			lineIndex = len(b.lines) - 1
		}
		curLine := b.lines[lineIndex]

		switch i {
		case 0: // Working with 1st line
			if len(b.copies) == 1 {
				if !isFirstStart && !isLastEnd {
					curLine.data = slices.Concat(dataFirst, data, dataSecond)
				} else {
					b.InsertLineWithData(lineIndex+shift, data)
				}
			} else {
				if isFirstStart && isLastEnd {
					b.InsertLineWithData(lineIndex+shift, data)
				} else {
					curLine.data = slices.Concat(dataFirst, data)
				}
			}
		case len(b.copies) - 1: // Working with last line
			if isLastEnd && isFirstStart {
				b.InsertLineWithData(lineIndex+shift, data)
			} else {
				savedData := slices.Concat(data, dataSecond)
				b.InsertLineWithData(lineIndex+shift, savedData)
			}
		default: //Just insert a new line
			b.InsertLineWithData(lineIndex+shift, data)
		}
	}

	b.fixOffset()
}
