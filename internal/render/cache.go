package render

import (
	"github.com/EnotInc/Bard/internal/enums"
)

// About |cachedLine|
// hash - hash of the raw text in line
// render - stores the rendered string, without Markdown symbols, but with ANSI characters
// diff - stores the difference in visible characters and original string length
// index - the line number
// mode - [render mode] that was used when cache is saved
type cachedLine struct {
	hash   uint32
	render string
	diff   int
	index  int
	mode   enums.Render
}

// About |cache|
// lines - list of [chchedLine]
// dirty - now it's just a `bool` value. Used to check if whole screen is needed to be randered again
type cache struct {
	lines map[int]*cachedLine
	dirty bool
}

func initCache() *cache {
	return &cache{lines: make(map[int]*cachedLine), dirty: false}
}

// About getCached()
// so here is where I get cached lines
// It returns 2 values, pointer to cached line and a boolean, wich use to see if line was cached before or not
// but I don't just return `l, ok := c.line[index]`. I also checks if render dirty or not (`ok && !c.dirty`), and if it is, render will think that line wasn't cached before, and it will render it again
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
