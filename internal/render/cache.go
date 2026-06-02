package render

import (
	render "github.com/EnotInc/Bard/internal/enums/render"
)

// hash - hash of the raw text in line
// render - stores the rendered string, without Markdown symbols, but with ANSI characters
// index - the line number
// mode - render mode that was used when cache is saved
type cachedLine struct {
	render string
	index  int
	mode   render.Render
	hash   uint32
	keep   bool
}

// lines - list of chchedLine
// dirty - now it's just a `bool` value. Used to check if whole screen is needed to be randered again
type cache struct {
	lines map[int]*cachedLine
	dirty bool
}

func initCache() *cache {
	return &cache{lines: make(map[int]*cachedLine), dirty: false}
}

// so here is where I get cached lines
// It returns 2 values, pointer to cached line and a boolean, wich use to see if line was cached before or not
// but I don't just return `l, ok := c.line[index]`. I also checks if render dirty or not (`ok && !c.dirty`), and if it is, render will think that line wasn't cached before, and it will render it again
func (c *cache) getCached(index int) (*cachedLine, bool) {
	l, ok := c.lines[index]
	return l, ok && !c.dirty
}

func (c *cache) cacheLine(h uint32, render string, index int, m render.Render, keep bool) {
	// If the line exists in the map, update it
	if l, ok := c.lines[index]; ok {
		l.hash = h
		l.render = render
		l.index = index
		l.mode = m
		l.keep = keep
	} else { // Otherwise, create a new one
		newLine := &cachedLine{}
		newLine.hash = h
		newLine.render = render
		newLine.index = index
		newLine.mode = m
		newLine.keep = keep

		c.lines[index] = newLine
	}
}

func (c *cache) purge() {
	c.dirty = false
	for k := range c.lines {
		delete(c.lines, k)
	}

}
