package buffer

type Buffer struct {
	Title      string
	pairs      []rune // paired brackets
	copies     []*copied
	Lines      []*Line
	Cursor     *cursor
	Visual     *cursor
	IsReadOnly bool
	IsMdFile   bool
}

func InitBuffer() []*Buffer {
	c := &cursor{line: 0, offset: 0}
	v := &cursor{line: 0, offset: 0}
	b := &Buffer{
		Cursor:     c,
		Visual:     v,
		pairs:      []rune{},
		IsReadOnly: false,
		IsMdFile:   false,
	}
	b.Lines = append(b.Lines, &Line{Data: []rune("")})
	var bfs []*Buffer
	bfs = append(bfs, b)
	return bfs
}

func (b *Buffer) EscapeToNormal() {
	b.pairs = []rune{}
	if b.Cursor.offset > 0 {
		b.Cursor.offset -= 1
	}
}

func (b *Buffer) ResetKeepOffset() {
	b.Cursor.keepOffset = b.Cursor.offset
}
