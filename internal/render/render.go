package render

import (
	"github.com/EnotInc/Bard/internal/enums"
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

func InitRender(w, h int) *Renderer {
	_c := initCache()
	r := &Renderer{c: _c, w: w, mode: enums.Markdown}
	r.md = md.NewRender(w)
	r.code = code.NewRender(w)
	return r
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

func (r *Renderer) Render(line []rune, lineIndex int, show bool, isCurrent bool, isFirst bool) (string, int) {
	lineHash := GetHash(&line)
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
		data, diff, mode = r.code.RenderCodeLine(line)
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
