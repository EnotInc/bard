package explorer

import (
	"strings"

	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/services"
)

const placeholder = "'/' to search"

const searchBarOfset = 1

func (ex *Explorer) buildSearchBar() string {
	var searchBar strings.Builder
	searchBar.WriteString(string(ascii.UnderLine))

	icon := services.SearchIcon()
	searchBar.WriteString(icon)

	if len(ex.search) == 0 {
		searchBar.WriteString(placeholder)

		amount := max(0, ex.w-len(placeholder))
		fill := strings.Repeat(" ", amount)
		searchBar.WriteString(fill)

	} else {
		searchBar.WriteString(string(ex.search))
		amount := max(0, ex.w-len(ex.search))
		fill := strings.Repeat(" ", amount)
		searchBar.WriteString(fill)
	}

	return searchBar.String()
}

func (ex *Explorer) beginSearch() {
	ex.action = searching
	ex.cursor.y = 0
	ex.yScroll = 0
}

func (ex *Explorer) handleSearch(key rune) {
	switch key {
	case keys.Esc, keys.Enter:
		ex.action = none
		ex.moveToTop()
		ex.fixCursor()
		ex.scroll()
	case keys.Backspace:
		if len(ex.search) > 0 {
			ex.search = ex.search[:len(ex.search)-1]
		}
	default:
		if services.IsLetterOrNumber(key) || key == '.' {
			ex.search = append(ex.search, key)
		}
	}
}
