package buffer

// About |Buffer|
// |Title| - name of the oppened file
// |pairs| - used like a stask to keep track of paired brackets
// |copies| - sores copied lines
// |Lines| - List of lines
// |Cursor| - read pos of cursor
// |Visual| - anchor point for real cursor pos, used to calculate selected area in visual and visual-line modes
// |IsReadOnly| - ised to check if oppened file could be change
// |IsMdFile| - if file extations is not '.md', bard uses default reader instead of markdown one
type Buffer struct {
	Title      string
	pairs      []rune
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
