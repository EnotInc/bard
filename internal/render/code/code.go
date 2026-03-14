package code

import (
	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/render/general"
)

type Render struct {
	w int
}

func NewRender(w int) *Render {
	return &Render{w: w}
}

func (r *Render) RenderCodeLine(line []rune) (string, enums.Render) {
	var mode enums.Render = enums.Code
	data := general.PaintString(ascii.YellowFg, string(line))
	if string(line) == "```" {
		mode = enums.Markdown
	}
	return data, mode
}
