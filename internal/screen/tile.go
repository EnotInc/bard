package screen

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
)

type object interface {
	DrawLineAt(index int) string
	Handle(key rune)
}

type tile struct {
	hash   map[int]uint32
	object object
	title  string
	w, h   int
	x, y   int
	border bool
}

func NewTile(o object, t string, w, h, x, y int, b bool) *tile {
	tile := &tile{
		object: o,
		title:  t,
		w:      w,
		h:      h,
		x:      x,
		y:      y,
		border: b,
	}
	return tile
}

func (t *tile) GetDiff() string {
	var diff strings.Builder

	for i := range t.w {
		var data strings.Builder
		if t.border && i == 0 {
			border := t.getBorder(true)
			data.WriteString(string(ascii.BorderCUL))
			data.WriteString(border)
			data.WriteString(string(ascii.BorderCUR))

			diff.WriteString(data.String())
			continue
		}
		if t.border && i == t.h {
			border := t.getBorder(false)
			data.WriteString(string(ascii.BorderCUL))
			data.WriteString(border)
			data.WriteString(string(ascii.BorderCUR))

			diff.WriteString(data.String())
			continue
		}

		line := t.object.DrawLineAt(i)
		ofset := 0
		if t.border {
			ofset = 1
		}
		trim := services.VisibleSubString(line, 0+ofset, t.w-ofset)

		// TODO: figure out if there is a better way to draw with borders
		if t.border {
			data.WriteString(string(ascii.BorderH))
		}
		data.WriteString(trim)
		if t.border {
			fmt.Fprintf(&data, "\033[%d;%dH", t.y+i, t.x)
			data.WriteString(string(ascii.BorderH))
		}

		curHash := services.GetHash(data.String())
		oldHash, ok := t.hash[i]
		if !ok || (ok && curHash != oldHash) { // add check for current line
			// TODO: figure out how to work with cursor
			pos := fmt.Sprintf("\033[%d;%dH\033[0K", t.x, t.y)
			diff.WriteString(pos)
		}
	}

	return diff.String()
}

func (t *tile) getBorder(withTitle bool) string {
	var border strings.Builder
	if withTitle {
		border.WriteString(string(ascii.BorderV))
		border.WriteString(t.title)
		border.WriteString(strings.Repeat(string(ascii.BorderH), t.w-1-len(t.title)))
	} else {
		border.WriteString(strings.Repeat(string(ascii.BorderH), t.w-2))
	}
	return border.String()
}
