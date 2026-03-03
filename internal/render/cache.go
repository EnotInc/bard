package render

import "slices"

type cachedLine struct {
	// raw is used for string comparison
	raw []rune

	// render stores the rendered string
	// without Markdown symbols, but with ANSI characters
	render string

	// diff stores the difference in visible
	// characters and original string length
	diff int

	// index - the line number
	index int
}

type cache struct {
	lines map[int]*cachedLine
}

func initCache() *cache {
	return &cache{lines: make(map[int]*cachedLine)}
}

func (b *cache) isCached(index int) bool {
	_, ok := b.lines[index]
	return ok
}

func (b *cache) getCached(index int) *cachedLine {
	l, _ := b.lines[index]
	return l
}

func (b *cache) cacheLine(raw []rune, render string, diff int, index int) {
	// If the line exists in the map, update it
	if l, ok := b.lines[index]; ok {
		l.raw = slices.Clone(raw)
		l.render = render
		l.diff = diff
		l.index = index
	} else { // Otherwise, create a new one
		newLine := &cachedLine{}
		newLine.raw = slices.Clone(raw)
		newLine.render = render
		newLine.diff = diff
		newLine.index = index

		b.lines[index] = newLine
	}
}
