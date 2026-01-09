package editor

type line struct {
	data []rune
	len  int
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
	b.lines = append(b.lines, &line{len: 0})
	return b
}

func (b *Buffer) H() {}
func (b *Buffer) J() {}
func (b *Buffer) K() {}
func (b *Buffer) L() {}

func (b *Buffer) InsertKey(key rune) {
	curLine := b.lines[b.cursor.line]
	curLine.data = append(curLine.data[:b.cursor.ofset], append([]rune{key}, curLine.data[b.cursor.ofset:]...)...)
	b.cursor.ofset += 1
}

func (b *Buffer) RemoveKey() {}

func (b *Buffer) InsertEmptyLine(lineShift int) {}

func (b *Buffer) InsertLine() {}

func (b *Buffer) DelAndMoveLine() {}

func (b *Buffer) RemoveLine() {}
