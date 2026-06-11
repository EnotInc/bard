package explorer

import (
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/screen"
)

type cursor struct {
	x, y int
}

func initCursor() *cursor {
	return &cursor{
		x: 0,
		y: 0,
	}
}

func (ex *Explorer) j() {
	ex.cursor.y += 1
	if ex.cursor.y > len(ex.entries)-1 {
		ex.cursor.y = len(ex.entries) - 1
	}

	ex.scroll()
}

func (ex *Explorer) k() {
	ex.cursor.y -= 1
	if ex.cursor.y < 0 {
		ex.cursor.y = 0
	}

	ex.scroll()
}

func (ex *Explorer) scroll() {
	if ex.yScroll > ex.cursor.y-enums.ScrollBorder {
		ex.yScroll = max(0, ex.cursor.y-enums.ScrollBorder)
	} else if ex.yScroll < ex.cursor.y+enums.ScrollBorder-screen.H() {
		ex.yScroll = ex.cursor.y + enums.ScrollBorder - screen.H()
	}

	ex.visible.y = ex.cursor.y - ex.yScroll - 1
}

func (ex *Explorer) fixCursor() {
	if ex.cursor.y < 0 {
		ex.cursor.y = 0
	}
	if ex.cursor.y > len(ex.entries)-1 {
		ex.cursor.y = len(ex.entries) - 1
	}
}
