package render

import (
	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/hash"
	code "github.com/EnotInc/Bard/internal/render/code"
	md "github.com/EnotInc/Bard/internal/render/markdown"
)

// About |Renderer|
// struct is used to work with different renders
// mode - current [render mode]
// c - [cache]
// md - markdown redner
// code - code redner
// w - screen width
type Renderer struct {
	mode enums.Render
	c    *cache
	md   *md.Render
	code *code.Render
	w    int
}

func InitRender(w, h int, theme *config.Theme) *Renderer {
	_c := initCache()
	r := &Renderer{c: _c, w: w, mode: enums.Markdown}
	r.md = md.NewRender(w, &theme.Markdown)
	r.code = code.NewRender(w, &theme.Code)
	return r
}

func (r *Renderer) MakeDirty() {
	r.c.dirty = true
}

func (r *Renderer) Resize(w int) {
	r.md.Resize(w)
	r.code.Resize(w)
}

func (r *Renderer) Reset() {
	r.mode = enums.Markdown
	r.md.Reset()
	r.code.Reset()
	r.c.dirty = false
}

// About Render()
// This func is used to decide which render to use, and should you ever call a [Code] or [Markdown] render, or this line was already rendered an cached
// First it calculates hash of current line, an if this line was cached, it does next:
// 1. If this is a first line in render (first on the screen, on top) and if this line was c `code` line - current render mode is become `code`. This needed to avoid situation, where code block is starts above the visiable screen, and render would thing that text on the screen is a Makrdown, and node a code block
// 2. If hash of the line is stored equal to cached line (and it is not dirty) - it returns old rendered line (with escape sequences). This way I can save some time on render line, which wasn't changed, and just return prev render of this line
//
// Afther that it comares current render mode, decide wich render to use
// If rednered line has change render mode (if '```' is found), render switches modes, and makes all lines bellow dirty
// And then - caches the result of the render
// Basically, I render only line with the cursor on it, and dirty lines
func (r *Renderer) Render(line []rune, lineIndex int, show bool, isCurrent bool, isFirst bool) (string, int) {
	lineHash := hash.GetHash(string(line))
	if !isCurrent {
		if l, ok := r.c.getCached(lineIndex); ok {
			if isFirst && l.mode == enums.Code {
				r.mode = enums.Code
			}
			if lineHash == l.hash && l.mode == r.mode {
				return l.render, l.diff
			}
		}
	}

	var data string
	var diff int
	var mode enums.Render

	switch r.mode {
	case enums.Markdown:
		data, diff, mode = r.md.RenderMarkdownLine(line, lineIndex, show)
	case enums.Code:
		data, diff, mode = r.code.RenderCodeLine(line, show)
	}

	// If mode has changed, lines below becomes dirty
	if r.mode != mode {
		r.mode = mode
		r.c.dirty = true
	}

	if !isCurrent {
		r.c.cacheLine(lineHash, data, diff, lineIndex, r.mode)
	}
	return data, diff
}

func (r *Renderer) PurgeCache() {
	r.c.purge()
}
