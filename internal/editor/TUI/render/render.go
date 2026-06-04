package render

import (
	code "github.com/EnotInc/Bard/internal/editor/TUI/render/code"
	md "github.com/EnotInc/Bard/internal/editor/TUI/render/markdown"
	render "github.com/EnotInc/Bard/internal/enums/render"
	"github.com/EnotInc/Bard/internal/services"
)

// struct is used to work with different renders
// mode - current render mode
// c - cache
// md - markdown redner
// code - code redner
// w - screen width
type Renderer struct {
	c    *cache
	md   *md.Render
	code *code.Render
	mode render.Render
	w    int
}

func InitRender(w, h int) *Renderer {
	_c := initCache()
	r := &Renderer{c: _c, w: w, mode: render.Markdown}
	r.md = md.NewRender(w)
	r.code = code.NewRender(w)
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
	r.mode = render.Markdown
	r.md.Reset()
	r.code.Reset()
	r.c.dirty = false
}

func (r *Renderer) ToggleRender() {
	switch r.mode {
	case render.Code:
		r.mode = render.Markdown
	case render.Markdown:
		r.mode = render.Code
	}
}

// This func is used to decide which render to use, and should you ever call a Code or Markdown render, or this line was already rendered an cached
// First it calculates hash of current line, an if this line was cached, it does next:
// 1. If this is a first line in render (first on the screen, on top) and if this line was c `code` line - current render mode is become `code`. This needed to avoid situation, where code block is starts above the visiable screen, and render would thing that text on the screen is a Makrdown, and node a code block
// 2. If hash of the line is stored equal to cached line (and it is not dirty) - it returns old rendered line (with escape sequences). This way I can save some time on render line, which wasn't changed, and just return prev render of this line
//
// Afther that it comares current render mode, decide wich render to use
// If rednered line has change render mode (if '```' is found), render switches modes, and makes all lines bellow dirty
// And then - caches the result of the render
// Basically, I render only line with the cursor on it, and dirty lines
func (r *Renderer) Render(line []rune, lineIndex int, show bool, isCurrent bool, isFirst bool, xOfset int) (string, bool) {
	//lineHash := services.GetHash(string(line))
	lineHash := services.GetHash(string(line))
	if !isCurrent {
		if l, ok := r.c.getCached(lineIndex); ok && !l.keep {
			if isFirst && l.mode == render.Code {
				r.mode = render.Code
			}
			if lineHash == l.hash && l.mode == r.mode {
				return l.render, l.keep
			}
		}
	}

	var data string
	var mode render.Render
	var keep = false

	switch r.mode {
	case render.Markdown:
		data, mode, keep = r.md.RenderMarkdownLine(line, lineIndex, show, xOfset)
	case render.Code:
		data, mode, keep = r.code.RenderCodeLine(line, show, xOfset)
	}

	// If mode has changed, lines below becomes dirty
	if r.mode != mode {
		r.mode = mode
		r.c.dirty = true
	}

	if !isCurrent {
		r.c.cacheLine(lineHash, data, lineIndex, r.mode, keep)
	}
	return data, keep
}

func (r *Renderer) PurgeCache() {
	r.c.purge()
}
