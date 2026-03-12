package render

import (
	"slices"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/render/general"
	md "github.com/EnotInc/Bard/internal/render/markdown"
)

type Renderer struct {
	mode enums.Render
	c    *cache
	md   *md.Render
	w    int
}

func InitRender(w, h int) *Renderer {
	_c := initCache()
	r := &Renderer{c: _c, w: w, mode: enums.Markdown}
	r.md = md.NewRender(w)
	return r
}

func (r *Renderer) Reset() {
	r.mode = enums.Markdown
	r.md.Reset()
	r.c.dirty = false
}

func (r *Renderer) Render(line []rune, lineIndex int, show bool) (string, int) {
	if !show {
		if l, ok := r.c.getCached(lineIndex); ok == true {
			if l.mode == r.mode && slices.Equal(l.raw, line) {
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
		if r.mode != mode {
			r.mode = mode
			r.c.dirty = true
		}
	case enums.Code:
		//TODO: move out to `code` folder. Expend with it's own lexer, tokens and other stuff
		data = general.PaintString(ascii.YellowFg, string(line))
		diff = 0

		if string(line) == "```" {
			r.mode = enums.Markdown
			r.c.dirty = true
		}
	}

	if !show {
		r.c.cacheLine(line, data, diff, lineIndex, r.mode)
	}
	return data, diff

}
