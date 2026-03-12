package render

import (
	"slices"

	"github.com/EnotInc/Bard/internal/ascii"
)

type mode string

const (
	markdown mode = "md"
	code     mode = "code"
)

type Renderer struct {
	curAttr string
	mode    mode
	c       *cache
	l       *lexer
	w       int
}

func InitRender(w, h int) *Renderer {
	_c := initCache()
	r := &Renderer{c: _c, w: w, mode: markdown}
	r.l = NewLexer()
	return r
}

func (r *Renderer) Reset() {
	// NOTE: ig it will groud soon or later
	r.curAttr = ascii.Reset.Str()
	r.mode = markdown
	r.l = NewLexer()
}

func (r *Renderer) Render(line []rune, lineIndex int, show bool) (string, int) {
	if !show && r.c.isCached(lineIndex) {
		l, ok := r.c.getCached(lineIndex)
		if ok && l.mode == r.mode && slices.Equal(l.raw, line) {
			return l.render, l.diff
		}
	}

	var data string
	var diff int

	switch r.mode {
	case markdown:
		data, diff = r.RenderMarkdownLine(line, lineIndex, show)
	case code:
		data = paintString(ascii.YellowFg, string(line))
		diff = 0

		if string(line) == "```" {
			r.mode = markdown
		}
	}

	if !show {
		r.c.cacheLine(line, data, diff, lineIndex, r.mode)
	}
	return data, diff

}
