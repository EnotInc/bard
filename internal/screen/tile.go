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
	GetCursor(withBorderl bool) (int, int)
	SetTitle() string
	PreDraw()
}

type tile struct {
	hash   map[int]uint32
	object object
	title  string
	w, h   int
	x, y   int
	border bool
}

func NewTile(o object, w, h, x, y int, b bool) (*tile, error) {
	if w+x > global.w || h+y > global.h {
		return nil, fmt.Errorf("Can't layout tile. Size is too big:\nw: %d\th: %d\nx: %d\ty: %d\n\n\rGlobal Screen:\nw: %d\th: %d",
			w, h, x, y, global.w, global.h)
	}
	tile := &tile{
		object: o,
		w:      w,
		h:      h,
		x:      x,
		y:      y,
		border: b,
	}
	return tile, nil
}

const termShift = 1

func (t *tile) GetDiff() string {
	var diff strings.Builder
	diff.WriteString(string(ascii.HideCursor))
	diff.WriteString(string(ascii.MoveToStart))

	statusLine := 1
	for i := range t.h - statusLine { // leaving one for status line
		var data strings.Builder
		if t.border && i == 0 {
			border := t.getBorder(true)
			data.WriteString(string(ascii.Reset))
			data.WriteString(string(ascii.BorderCUL))
			data.WriteString(border)
			data.WriteString(string(ascii.BorderCUR))

			pos := fmt.Sprintf("\033[%d;%dH", t.y+termShift, t.x+termShift)
			diff.WriteString(pos)
			diff.WriteString(data.String())
			continue
		}
		if t.border && i == t.h-1-statusLine {
			border := t.getBorder(false)
			pos := fmt.Sprintf("\033[%d;%dH", t.y+t.h-statusLine, t.x+termShift)
			diff.WriteString(pos)

			data.WriteString(string(ascii.Reset))
			data.WriteString(string(ascii.BorderCDL))
			data.WriteString(border)
			data.WriteString(string(ascii.BorderCDR))

			diff.WriteString(data.String())
			break
		}

		ofset := 0
		if t.border {
			ofset = 1
		}
		line := t.object.DrawLineAt(i - ofset)

		trim := services.VisibleSubString(line, 0, t.w-ofset*2)
		curHash := services.GetHash(trim)
		oldHash, ok := t.hash[i]

		if !ok || (ok && curHash != oldHash) { // add check for current line
			pos := fmt.Sprintf("\033[%d;%dH\033[0K", t.y+i+termShift, t.x+termShift)
			diff.WriteString(pos)
			if t.border {
				data.WriteString(string(ascii.Reset))
				data.WriteString(string(ascii.BorderV))
			}
			data.WriteString(trim)
			if t.border {
				fmt.Fprintf(&data, "\033[%d;%dH", t.y+i+termShift, t.x+t.w)
				data.WriteString(string(ascii.Reset))
				data.WriteString(string(ascii.BorderV))
			}

			diff.WriteString(data.String())
		}
	}

	return diff.String()
}

func (t *tile) getBorder(withTitle bool) string {
	var border strings.Builder
	if withTitle {
		t.title = t.object.SetTitle()
		visible := services.CountClear(t.title, 0, len(t.title))
		if visible >= t.w-2-termShift {
			t.title = services.VisibleSubString(t.title, 0, t.w-2-termShift-1)
			visible = services.CountClear(t.title, 0, len(t.title)-1)
		}
		amount := max(t.w-2-visible-termShift, 0)

		border.WriteString(string(ascii.BorderH))
		border.WriteString(t.title)
		border.WriteString(string(ascii.Reset))
		border.WriteString(strings.Repeat(string(ascii.BorderH), amount))
	} else {
		border.WriteString(strings.Repeat(string(ascii.BorderH), t.w-2))
	}
	return border.String()
}
