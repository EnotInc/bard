package editor

import (
	"slices"
)

const (
	head = true
	tail = false
)

type line struct {
	data []rune
}

type copyed struct {
	data  []rune
	start int
	end   int
}

type cursor struct {
	line  int
	ofset int
}

type Buffer struct {
	copyes []*copyed
	lines  []*line
	cursor *cursor
	visual *cursor
}

func InitBuffer() *Buffer {
	c := &cursor{line: 0, ofset: 0}
	v := &cursor{line: 0, ofset: 0}
	b := &Buffer{
		cursor: c,
		visual: v,
	}
	b.lines = append(b.lines, &line{data: []rune("")})
	return b
}

func (b *Buffer) H(amount int) {
	if b.cursor.ofset-amount > 0 {
		b.cursor.ofset -= amount
	} else {
		b.cursor.ofset = 0
	}
}

func (b *Buffer) J(amount int) {
	if b.cursor.line+amount < len(b.lines)-1 {
		b.cursor.line += amount
	} else {
		b.cursor.line = len(b.lines) - 1
	}
	if b.cursor.ofset > len(b.lines[b.cursor.line].data)-1 {
		b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
		if b.cursor.ofset < 0 {
			b.cursor.ofset = 0
		}
	}
}

func (b *Buffer) K(amount int) {
	if b.cursor.line-amount > 0 {
		b.cursor.line -= amount
	} else {
		b.cursor.line = 0
	}
	if b.cursor.ofset > len(b.lines[b.cursor.line].data)-1 {
		b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
		if b.cursor.ofset < 0 {
			b.cursor.ofset = 0
		}
	}
}

func (b *Buffer) L(amount int) {
	if b.cursor.ofset+amount < len(b.lines[b.cursor.line].data)-1 {
		b.cursor.ofset += amount
	} else {
		b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
		if b.cursor.ofset < 0 {
			b.cursor.ofset = 0
		}
	}
}

func (b *Buffer) InsertKey(key rune) {
	curLine := b.lines[b.cursor.line]
	curLine.data = append(curLine.data[:b.cursor.ofset], append([]rune{key}, curLine.data[b.cursor.ofset:]...)...)
	b.cursor.ofset += 1
}

func (b *Buffer) RemoveKey(keyShift int) {
	if b.cursor.ofset > 0 {
		curLine := b.lines[b.cursor.line]
		index := keyShift + b.cursor.ofset
		curLine.data = slices.Delete(curLine.data, index-1, index)
		b.cursor.ofset -= 1
	} else {
		b.DelAndMoveLine()
	}
}

func (b *Buffer) Delkey() {
	if len(b.lines[b.cursor.line].data) > 0 {
		curLine := b.lines[b.cursor.line]
		index := b.cursor.ofset
		curLine.data = slices.Delete(curLine.data, index, index+1)
	}
}

func (b *Buffer) DelRangeKey(startLine int, endLine int, startOfset int, endOfset int, i int) {
	if i < 0 || i >= len(b.lines) {
		return
	}

	curLine := b.lines[i]
	if len(curLine.data) == 0 {
		return
	}

	if startOfset < 0 {
		startOfset = 0
	}
	if endOfset > len(curLine.data) {
		endOfset = len(curLine.data)
	}

	curLine.data = slices.Delete(curLine.data, startOfset, endOfset)

	if len(curLine.data) == 0 {
		b.RemoveLineAt(i)
		if i == endLine {
			b.DelAndMoveLineAt(startLine, endLine, endOfset)
		}
	}
}

func (b *Buffer) InsertEmptyLine(lineShift int) {
	index := b.cursor.line + lineShift
	newLine := make([]*line, 0)
	newLine = append(newLine, &line{data: []rune("")})
	b.lines = append(b.lines[:index], append(newLine, b.lines[index:]...)...)
}

func (b *Buffer) InsertLine() {
	index := b.cursor.line + 1
	shiftData := b.lines[b.cursor.line].data[b.cursor.ofset:]
	b.lines[b.cursor.line].data = b.lines[b.cursor.line].data[:b.cursor.ofset]

	newLine := line{}
	b.lines = append(b.lines[:index], append([]*line{&newLine}, b.lines[index:]...)...)
	b.cursor.line += 1
	b.cursor.ofset = 0

	b.lines[b.cursor.line].data = append(b.lines[b.cursor.line].data, shiftData...)
}

func (b *Buffer) DelAndMoveLine() {
	if b.cursor.line > 0 {
		shiftData := b.lines[b.cursor.line].data[b.cursor.ofset:]
		b.RemoveLine()
		b.cursor.ofset = len(b.lines[b.cursor.line].data)
		b.lines[b.cursor.line].data = append(b.lines[b.cursor.line].data, shiftData...)
	}
}

func (b *Buffer) DelAndMoveLineAt(startLine int, endLine int, endOfset int) {
	shiftData := b.lines[endLine].data[endOfset:]
	b.RemoveLineAt(endLine)
	b.cursor.ofset = len(b.lines[startLine].data)
	b.lines[startLine].data = append(b.lines[startLine].data, shiftData...)
}

func (b *Buffer) RemoveLine() {
	if len(b.lines) == 1 {
		b.ClearLine()
		return
	}
	b.lines = slices.Delete(b.lines, b.cursor.line, b.cursor.line+1)
	if b.cursor.line >= len(b.lines) {
		b.K(1)
	}
}

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

func (b *Buffer) ClearLine() {
	b.lines[b.cursor.line].data = []rune{}
}

func (b *Buffer) moveToFirst() {
	for i := range len(b.lines[b.cursor.line].data) {
		if b.lines[b.cursor.line].data[i] != ' ' {
			b.cursor.ofset = i
			break
		}
	}
}

func (b *Buffer) copyLine(l *line, startOfset int, endOfset int) *copyed {
	if l == nil {
		return &copyed{data: []rune(""), start: 0, end: 0}
	}

	if startOfset < 0 {
		startOfset = 0
	}
	if endOfset > len(l.data) {
		endOfset = len(l.data)
	}
	if startOfset > endOfset {
		startOfset, endOfset = endOfset, startOfset
	}

	data := append([]rune(nil), l.data[startOfset:endOfset]...)
	return &copyed{data: data, end: endOfset, start: startOfset}
}

func (b *Buffer) copySelected() {
	b.copyes = []*copyed{}

	startOfset := b.visual.ofset
	startLine := b.visual.line
	endOfset := b.cursor.ofset
	endLine := b.cursor.line

	if startLine > endLine || (startLine == endLine && startOfset > endOfset) {
		startLine, endLine = endLine, startLine
		startOfset, endOfset = endOfset, startOfset
	}

	lineCount := 0
	lineSelected := endLine - startLine

	for i := startLine; i <= endLine; i++ {
		curLineStart := startOfset
		curLineEnd := endOfset + 1

		if i != startLine && lineCount != 0 {
			curLineStart = 0
		}
		if i != endLine && lineCount != lineSelected {
			curLineEnd = len(b.lines[i].data)
		}

		line := b.copyLine(b.lines[i], curLineStart, curLineEnd)
		b.copyes = append(b.copyes, line)
	}

	b.cursor.line = startLine
	b.cursor.ofset = startOfset
}

func (b *Buffer) deleteSelected() {
	startOfset := b.visual.ofset
	startLine := b.visual.line
	endOfset := b.cursor.ofset
	endLine := b.cursor.line

	if startLine > endLine || (startLine == endLine && startOfset > endOfset) {
		startLine, endLine = endLine, startLine
		startOfset, endOfset = endOfset, startOfset
	}
}
