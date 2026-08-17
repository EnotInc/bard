package screen

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
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
	cfg := config.GetConfig()
	theme := config.GetTheme().General

	for i := range p.h {

		if i == 0 && cfg.ShowBorder {
			x, y := p.calcPos()
			fmt.Fprintf(&diff, "\033[%d;%dH", y, x)
			diff.WriteString(theme.BottomBar)
			diff.WriteString(theme.SelectedTile)
			diff.WriteString(ascii.BorderCUL)
			diff.WriteString(strings.Repeat(ascii.BorderH, p.w-2))
			diff.WriteString(ascii.BorderCUR)
			diff.WriteString(ascii.Reset.Str())
			continue
		}

		if i == p.h-1 && cfg.ShowBorder {
			x, y := p.calcPos()
			fmt.Fprintf(&diff, "\033[%d;%dH", y+p.h-1, x)
			diff.WriteString(theme.BottomBar)
			diff.WriteString(theme.SelectedTile)
			diff.WriteString(ascii.BorderCDL)
			diff.WriteString(strings.Repeat(ascii.BorderH, p.w-2))
			diff.WriteString(ascii.BorderCDR)
			diff.WriteString(ascii.Reset.Str())
			break
		}

		offset := 0
		if cfg.ShowBorder {
			offset = 1
		}

		_l := p.object.DrawLineAt(i - offset)
		trim := services.VisibleSubString(_l, 0, p.w-offset*3)

		curHash := services.GetHash(trim)
		oldHash, ok := p.hash[i]

		if !ok || (ok && curHash != oldHash) {
			x, y := p.calcPos()
			pos := fmt.Sprintf("\033[%d;%dH", i+y, x)
			diff.WriteString(pos)
			if cfg.ShowBorder {
				diff.WriteString(theme.BottomBar)
				diff.WriteString(theme.SelectedTile)
				diff.WriteString(ascii.BorderV)
				diff.WriteString(ascii.Reset.Str())
			}

			diff.WriteString(trim)
			diff.WriteString(ascii.Reset.Str())

			if cfg.ShowBorder {
				fmt.Fprintf(&diff, "\033[%d;%dH", i+y, p.w+x-1)
				diff.WriteString(theme.BottomBar)
				diff.WriteString(theme.SelectedTile)
				diff.WriteString(ascii.BorderV)
				diff.WriteString(ascii.Reset.Str())
			}
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
