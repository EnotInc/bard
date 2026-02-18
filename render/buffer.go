package render

type cachedLine struct {
	raw    []rune
	render string
	diff   int
	index  int
}

type buffer struct {
	lines map[int]*cachedLine
}

func initBuffer() *buffer {
	return &buffer{lines: make(map[int]*cachedLine)}
}

func (b *buffer) isCached(index int) bool {
	_, ok := b.lines[index]
	return ok
}

func (b *buffer) getCached(index int) *cachedLine {
	l, _ := b.lines[index]
	return l
}

func (b *buffer) cacheLine(raw []rune, render string, diff int, index int) {
	if _, ok := b.lines[index]; ok {
		l := b.lines[index]
		l.raw = raw
		l.render = render
		l.diff = diff
		l.index = index
	} else {
		newLine := &cachedLine{}
		newLine.raw = raw
		newLine.render = render
		newLine.diff = diff
		newLine.index = index

		b.lines[index] = newLine
	}
}

