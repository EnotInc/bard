package render

import (
	"github.com/EnotInc/Bard/internal/enums"
	code "github.com/EnotInc/Bard/internal/render/code"
	md "github.com/EnotInc/Bard/internal/render/markdown"
)

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

func (r *Renderer) Reset() {
	r.mode = enums.Markdown
	r.md.Reset()
	r.c.dirty = false
}

func (r *Renderer) Render(line []rune, lineIndex int, show bool) (string, int) {
	lineHash := GetHash(&line)
	if !show {
		if l, ok := r.c.getCached(lineIndex); ok == true {
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
		//TODO: move out to `code` folder. Expend with it's own lexer, tokens and other stuff
		data, mode = r.code.RenderCodeLine(line)
		diff = 0

	}

	// If mode has changed, lines below becomes dirty
	if r.mode != mode {
		r.mode = mode
		r.c.dirty = true
	}

	if !show {
		r.c.cacheLine(lineHash, data, diff, lineIndex, r.mode)
	}
	return data, diff
}
