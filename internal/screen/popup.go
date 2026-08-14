package screen

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/services"
)

type popup struct {
	hash   map[int]uint32
	object object
	title  string
	w, h   int
}

func (p *popup) Draw() {
	var diff strings.Builder

	for i := range p.h {
		_l := p.object.DrawLineAt(i)
		trim := services.VisibleSubString(_l, 0, p.w)

		curHash := services.GetHash(trim)
		oldHash, ok := p.hash[i]

		if !ok || (ok && curHash != oldHash) {
			x, y := p.calcPos()
			pos := fmt.Sprintf("\033[%d;%dH", i+y, x)
			diff.WriteString(pos)
			diff.WriteString(trim)
		}
	}

	fmt.Print(diff.String())
}

const width int = 48

func calcSize() (w, h int) {
	_w := width
	_h := float32(global.h) * 0.8
	w, h = _w, int(_h)
	return w, h
}

func (p *popup) calcPos() (x, y int) {
	x = (global.w - p.w) / 2
	y = (global.h - p.h) / 2
	return x, y
}

func (p *popup) Resize() {
	_w, _h := calcSize()

	p.object.Resize(_w, _h)
}

func NewPopup(o object) *popup {
	_w, _h := calcSize()
	o.Resize(_w, _h)

	return &popup{
		object: o,
		w:      _w,
		h:      _h,
	}
}
