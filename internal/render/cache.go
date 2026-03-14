package render

import (
	"github.com/EnotInc/Bard/internal/enums"
)

type cachedLine struct {
	// Now I store hash instead the whole line
	hash uint32

	// render stores the rendered string
	// without Markdown symbols, but with ANSI characters
	render string

	// diff stores the difference in visible
	// characters and original string length
	diff int

	// index - the line number
	index int

	// mode that was used when cache is saved
	mode enums.Render
}

type cache struct {
	lines map[int]*cachedLine
	dirty bool
}

func initCache() *cache {
	return &cache{lines: make(map[int]*cachedLine), dirty: false}
}

func (c *cache) getCached(index int) (*cachedLine, bool) {
	l, ok := c.lines[index]
	return l, ok && !c.dirty
}

func (c *cache) cacheLine(h uint32, render string, diff int, index int, m enums.Render) {
	// If the line exists in the map, update it
	//var foo uint32
	if l, ok := c.lines[index]; ok {
		l.hash = h
		l.render = render
		l.diff = diff
		l.index = index
		l.mode = m
	} else { // Otherwise, create a new one
		newLine := &cachedLine{}
		newLine.hash = h
		newLine.render = render
		newLine.diff = diff
		newLine.index = index
		newLine.mode = m

		c.lines[index] = newLine
	}
}
