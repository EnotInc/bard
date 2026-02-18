package render

type cashedLine struct {
	raw    []rune
	render string
	diff   int
	index  int
}

type buffer struct {
	lines map[int]*cashedLine
}

func initBuffer() *buffer {
	return &buffer{lines: make(map[int]*cashedLine)}
}

func (b *buffer) isCashed(index int) bool {
	_, ok := b.lines[index]
	return ok
}

func (b *buffer) getCashed(index int) *cashedLine {
	l, _ := b.lines[index]
	return l
}

func (b *buffer) casheLine(raw []rune, render string, diff int, index int) {
	if _, ok := b.lines[index]; ok {
		l := b.lines[index]
		l.raw = raw
		l.render = render
		l.diff = diff
		l.index = index
	} else {
		new := &cashedLine{}
		new.raw = raw
		new.render = render
		new.diff = diff
		new.index = index

		b.lines[index] = new
	}
}
