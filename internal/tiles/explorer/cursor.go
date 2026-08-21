package explorer

import (
	"slices"

	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/screen"
)

type cursor struct {
	x, y int
}

func initCursor() *cursor {
	return &cursor{
		x: 0,
		y: searchBarOfset,
	}
}

func (ex *Explorer) j() {
	ex.cursor.y += 1
	if ex.cursor.y > len(ex.entries) {
		ex.cursor.y = len(ex.entries)
	}

	ex.scroll()
}

func (ex *Explorer) k() {
	ex.cursor.y -= 1
	if ex.cursor.y < searchBarOfset {
		ex.cursor.y = searchBarOfset
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
	if ex.cursor.y < searchBarOfset {
		ex.cursor.y = searchBarOfset
	}
	if ex.cursor.y > len(ex.entries) {
		ex.cursor.y = len(ex.entries)
	}
}

func (ex *Explorer) moveToTop() {
	ex.cursor.y = searchBarOfset
	if !slices.Equal(ex.path, screen.Root()) && len(ex.entries) > 1 {
		ex.cursor.y += 1
	}
}

func (ex *Explorer) moveToBottom() {
	ex.cursor.y = max(len(ex.entries), 0)
} // ex.cursor.y = max(len(ex.enties)-1+serchBarOfset, 0)
