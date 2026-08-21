package themes

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
	"github.com/EnotInc/Bard/theme"
)

const placeholder = "'/' to search"

const searchOffset = 2

func (t *Themes) buildSearchBar() string {
	var searchBar strings.Builder

	searchBar.WriteString(ascii.UnderLine.Str())
	searchBar.WriteString(theme.GetTheme().General.BottomBar)
	si := config.GetConfig().ShowIcons
	icon := services.SearchIcon(si)
	searchBar.WriteString(" ")
	searchBar.WriteString(icon)

	if len(t.search) == 0 && t.action == none {
		searchBar.WriteString(placeholder)
		amount := max(0, t.w-len(placeholder)-1-searchOffset)
		fill := strings.Repeat(" ", amount)
		searchBar.WriteString(fill)

	} else {
		searchBar.WriteString(string(t.search))
		amount := max(0, t.w-len(t.search)-1-searchOffset)
		fill := strings.Repeat(" ", amount)
		searchBar.WriteString(fill)
	}

	searchBar.WriteString(ascii.Reset.Str())
	return searchBar.String()
}
