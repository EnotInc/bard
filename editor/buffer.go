package editor

import (
	"slices"
)

type line struct {
	data []rune
}

type copyed struct {
	data    []rune
	isEnd   bool
	isStart bool
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

func (b *Buffer) fixOfset() {
	if b.cursor.ofset > len(b.lines[b.cursor.line].data)-1 {
		b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
	}
	if b.cursor.ofset < 0 {
		b.cursor.ofset = 0
	}
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
	b.fixOfset()
}

func (b *Buffer) K(amount int) {
	if b.cursor.line-amount > 0 {
		b.cursor.line -= amount
	} else {
		b.cursor.line = 0
	}
	b.fixOfset()
}

func (b *Buffer) L(amount int) {
	if b.cursor.ofset+amount < len(b.lines[b.cursor.line].data)-1 {
		b.cursor.ofset += amount
	} else {
		b.cursor.ofset = len(b.lines[b.cursor.line].data) - 1
		b.fixOfset()
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
		ch := curLine.data[index]
		b.copyes = append([]*copyed{}, &copyed{data: []rune{ch}, isStart: false, isEnd: false})
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
	_isStart := startOfset == 0
	_isEnd := endOfset == len(l.data)

	if startOfset < 0 {
		startOfset = 0
	}
	if endOfset > len(l.data)-1 {
		endOfset = len(l.data)
	}
	if endOfset < 0 {
		endOfset = 0
	}
	if startOfset > endOfset {
		startOfset, endOfset = endOfset, startOfset
	}
	var newData []rune
	// if startOfset == endOfset || len(l.data) == 0 {
	// 	newData = []rune("")
	//} else {
	newData = l.data[startOfset:endOfset]
	if len(newData) == 0 {
		newData = []rune("")
	}
	//}
	data := append([]rune{}, newData...)
	return &copyed{data: data, isEnd: _isEnd, isStart: _isStart}
}

func (b *Buffer) copySelected(isDelete bool, isVisualLine bool) {
	b.copyes = []*copyed{}

	startOfset := b.visual.ofset
	startLine := b.visual.line
	endOfset := b.cursor.ofset
	endLine := b.cursor.line

	if startLine > endLine || (startLine == endLine && startOfset > endOfset) {
		startLine, endLine = endLine, startLine
		startOfset, endOfset = endOfset, startOfset
	}

	var tempLine []rune

	lineCount := 0
	lineSelected := endLine - startLine
	for i := startLine; i <= endLine; {
		curOfsetStart := 0
		if lineCount == 0 {
			curOfsetStart = startOfset
		}

		curOfsetEnd := max(len(b.lines[i].data)-1, 0)
		if i == endLine && !isVisualLine {
			curOfsetEnd = endOfset
		}
		if len(b.lines[i].data) > 0 {
			curOfsetEnd++
		}

		line := b.copyLine(b.lines[i], curOfsetStart, curOfsetEnd)
		b.copyes = append(b.copyes, line)

		if isDelete {
			if lineCount == 0 {
				tempLine = b.lines[i].data[:curOfsetStart]
			}
			if lineCount == lineSelected {
				tempLine = append(tempLine, b.lines[i].data[curOfsetEnd:]...)
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

	b.cursor.ofset = startOfset
	b.fixOfset()
}

func (b *Buffer) paste(shift int) {
	initialOfset := b.cursor.ofset + shift
	if len(b.lines[b.cursor.line].data) == 0 || b.cursor.ofset < 0 {
		initialOfset = 0
	}

	dataFirst := append([]rune(nil), b.lines[b.cursor.line].data[:initialOfset]...)
	dataSecond := append([]rune(nil), b.lines[b.cursor.line].data[initialOfset:]...)

	isFisrtStart := b.copyes[0].isStart
	isLastEnd := b.copyes[len(b.copyes)-1].isEnd

	lineIndex := b.cursor.line
	//ofsetIndex := len(dataSecond)
	for i, l := range b.copyes {
		lineIndex = b.cursor.line + i
		if lineIndex >= len(b.lines) {
			lineIndex = len(b.lines) - 1
		}
		curLine := b.lines[lineIndex]
		switch i {
		case 0:
			if len(b.copyes) == 1 {
				if !isFisrtStart && !isLastEnd {
					curLine.data = slices.Concat(dataFirst, l.data, dataSecond)
				} else {
					b.InsertLineWithData(lineIndex+shift, l.data)
				}
				//ofsetIndex = len(dataFirst) + len(l.data)
			} else {
				if isFisrtStart && isLastEnd {
					b.InsertLineWithData(lineIndex+shift, l.data)
				} else {
					curLine.data = slices.Concat(dataFirst, l.data)
				}
			}
		case len(b.copyes) - 1:
			if isLastEnd && isFisrtStart {
				b.InsertLineWithData(lineIndex+shift, l.data)
			} else {
				savedData := slices.Concat(l.data, dataSecond)
				b.InsertLineWithData(lineIndex+shift, savedData)
			}
		default:
			b.InsertLineWithData(lineIndex+shift, l.data)
		}
	}

	//b.cursor.line = min(lineIndex, len(b.lines)-1)
	//b.cursor.ofset = ofsetIndex
	b.fixOfset()
}
