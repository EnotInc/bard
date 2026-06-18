package screen

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
)

type object interface {
	DrawLineAt(index int) string
	Handle(key rune)
	GetCursor(withBorder bool) (int, int)
	SetTitle() string
	Resize(w, h int)
	PreDraw()
}

type tile struct {
	hash    map[int]uint32
	object  object
	title   string
	spacing float32
	w, h    int
}

func calcWithSpacing(spacing float32) (int, error) {
	if 0.0 > spacing || spacing > 1.0 {
		return 0, fmt.Errorf("spacing can't be less then 0, or greater then 1\nCurrent value: %f", spacing)
	}

	_w := float32(global.w) * spacing
	return int(_w), nil
}

func NewTile(o object, spacing float32) (*tile, error) {
	_h := global.h
	_w, err := calcWithSpacing(spacing)
	if err != nil {
		return nil, err
	}

	tile := &tile{
		object:  o,
		spacing: spacing,
		w:       _w,
		h:       _h,
	}
	return tile, nil
}

// move to enums?
const termShift = 1

func (t *tile) GetDiff(tileOfset int, isFocused bool) string {
	var diff strings.Builder
	diff.WriteString(string(ascii.HideCursor))
	diff.WriteString(string(ascii.MoveToStart))

	border := config.GetConfig().ShowBorder

	statusLine := 1
	for i := range t.h - statusLine { // leaving one for status line
		c := t.getColor(isFocused)
		var data strings.Builder
		if border && i == 0 {
			border := t.getBorder(true, isFocused)
			data.WriteString(string(ascii.Reset))
			data.WriteString(c)
			data.WriteString(string(ascii.BorderCUL))
			data.WriteString(border)
			data.WriteString(c)
			data.WriteString(string(ascii.BorderCUR))

			pos := fmt.Sprintf("\033[%d;%dH", termShift, tileOfset+termShift)
			diff.WriteString(pos)
			diff.WriteString(data.String())
			continue
		}
		if border && i == t.h-1-statusLine {
			border := t.getBorder(false, isFocused)
			pos := fmt.Sprintf("\033[%d;%dH", t.h-statusLine, tileOfset+termShift)
			diff.WriteString(pos)

			data.WriteString(string(ascii.Reset))
			data.WriteString(c)
			data.WriteString(string(ascii.BorderCDL))
			data.WriteString(border)
			data.WriteString(c)
			data.WriteString(string(ascii.BorderCDR))

			diff.WriteString(data.String())
			break
		}

		offset := 0
		if border {
			offset = 1 // NOTE: border ofset
		}
		line := t.object.DrawLineAt(i - offset)

		// NOTE: magic 2
		trim := services.VisibleSubString(line, 0, t.w-2)
		if !border { // used to make borderlett tiles more readable
			trim = fmt.Sprintf(" %s", trim)
		}

		curHash := services.GetHash(trim)
		oldHash, ok := t.hash[i]

		if !ok || (ok && curHash != oldHash) { // add check for current line
			pos := fmt.Sprintf("\033[%d;%dH\033[0K", i+termShift, tileOfset+termShift)
			diff.WriteString(pos)
			if border {
				data.WriteString(c)
				data.WriteString(string(ascii.BorderV))
				data.WriteString(string(ascii.Reset))
			}
			data.WriteString(trim)
			if border {
				fmt.Fprintf(&data, "\033[%d;%dH", i+termShift, tileOfset+t.w)
				data.WriteString(string(ascii.Reset))
				data.WriteString(c)
				data.WriteString(string(ascii.BorderV))
			}
			data.WriteString(string(ascii.Reset))

			diff.WriteString(data.String())
		}
	}

	return diff.String()
}

func (t *tile) getBorder(withTitle bool, isFocused bool) string {
	var border strings.Builder
	c := t.getColor(isFocused)
	if withTitle {
		t.title = t.object.SetTitle()
		visible := services.CountClear(t.title, 0, len(t.title))
		if visible >= t.w-2-termShift {
			t.title = services.VisibleSubString(t.title, 0, t.w-2-termShift-1)
			visible = services.CountClear(t.title, 0, len(t.title)-1)
		}
		amount := max(t.w-2-visible-termShift, 0)

		border.WriteString(c)
		border.WriteString(string(ascii.BorderH))
		border.WriteString(ascii.Reset.Str())
		border.WriteString(t.title)
		border.WriteString(string(ascii.Reset))
		border.WriteString(c)
		border.WriteString(strings.Repeat(string(ascii.BorderH), amount))
	} else {
		border.WriteString(c)
		border.WriteString(strings.Repeat(string(ascii.BorderH), t.w-2))
	}
	return border.String()
}

func (t *tile) getColor(focused bool) string {
	if focused {
		theme := config.GetTheme()
		return theme.General.SelectedTile
	} else {
		return ""
	}
}
