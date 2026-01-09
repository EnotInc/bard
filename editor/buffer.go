package editor

import (
	"slices"
)

type line struct {
	data []rune
}

type cursor struct {
	line  int
	ofset int
}

type Buffer struct {
	lines  []*line
	cursor *cursor
}

func InitBuffer() *Buffer {
	c := &cursor{line: 0, ofset: 0}
	b := &Buffer{
		cursor: c,
	}
	b.lines = append(b.lines, &line{data: []rune("")})
	return b
}

// NOTE: maybe this is not the best solution. But it's the only one that I came up with
// NOTE: this func is useful in case, when u jump up or down the line, and if it's empty cursor will be set to -1
func (b *Buffer) fixCursor() {
	if b.cursor.ofset < 0 {
		b.cursor.ofset = 0
	}
}

func (b *Buffer) H() {
	if b.cursor.ofset > 0 {
		b.cursor.ofset -= 1
	} else if b.cursor.line > 0 {
		b.cursor.line -= 1
		b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
	}
}

func (b *Buffer) J() {
	if b.cursor.line < len(b.lines)-1 {
		b.cursor.line += 1
		if b.cursor.ofset > len(b.lines[b.cursor.line].data)-1 {
			b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
			b.fixCursor()
		}
	}
}

func (b *Buffer) K() {
	if b.cursor.line > 0 {
		b.cursor.line -= 1
		if b.cursor.ofset > len(b.lines[b.cursor.line].data)-1 {
			b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
			b.fixCursor()
		}
	}
}

func (b *Buffer) L() {
	if b.cursor.ofset < len(b.lines[b.cursor.line].data)-1 { //NOTE: -1 here is used to not get out of bounds
		b.cursor.ofset += 1
	} else if b.cursor.line < len(b.lines)-1 {
		b.cursor.line += 1
		b.cursor.ofset = 0
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
	curLine := b.lines[b.cursor.line]
	index := b.cursor.ofset
	curLine.data = slices.Delete(curLine.data, index, index+1)
	if b.cursor.ofset == len(b.lines[b.cursor.line].data) {
		b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
		b.fixCursor()
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
		b.cursor.line -= 1
		b.cursor.ofset = len(b.lines[b.cursor.line].data)
		b.lines[b.cursor.line].data = append(b.lines[b.cursor.line].data, shiftData...)
	}
}

func (b *Buffer) RemoveLine() {
	b.lines = slices.Delete(b.lines, b.cursor.line, b.cursor.line+1)
}
