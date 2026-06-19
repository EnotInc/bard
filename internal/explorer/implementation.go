package explorer

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	cursorType "github.com/EnotInc/Bard/internal/enums/cursor"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"
)

func (ex *Explorer) DrawLineAt(index int) string {
	if index+ex.yScroll > len(ex.entries) {
		return ""
	}

	cfg := config.GetConfig()
	ofset := 2
	if cfg.ShowBorder {
		ofset *= 2
	}

	if index == 0 {
		return ex.buildSearchBar()
	}

	if ex.yScroll > 0 && index == 1 {
		return ascii.ArrowUp.Str()
	} else if index == ex.h-ofset {
		return ascii.ArrowDown.Str()
	}

	entry := ex.entries[index+ex.yScroll-searchBarOfset]
	var icon string
	if entry.isDir {
		icon = services.GetDirIcon(string(entry.name))
	} else {
		icon = services.GetFileIcon(string(entry.name))
	}

	if ex.action == deleting && index == ex.visible.y+1 {
		red := "\033[31m"
		green := "\033[32m"
		icon = fmt.Sprintf(" %sy%s/%sn%s: %s%s", red, ascii.ResetFg, green, ascii.Reset, ascii.Stricked, icon)
	}

	e := fmt.Sprintf("%s%s", icon, string(entry.name))
	e = services.VisibleSubString(e, 0, ex.w)
	return e
}

func (ex *Explorer) Handle(key rune) {
	switch ex.action {
	case creating:
		ex.handleCreate(key)
		return
	case changing:
		ex.handleChanges(key)
		return
	case deleting:
		ex.handleDeletion(key)
		return
	case searching:
		ex.handleSearch(key)
		return
	}

	switch key {
	case keys.Esc:
		if len(ex.search) > 0 {
			ex.search = []rune{}
			ex.update = true
		} else {
			screen.SetFocus(0)
		}
	case keys.Enter:
		ex.openEntryWithCallback()
		ex.moveToTop()
	case 'o':
		ex.beginCreation()
	case 's':
		ex.beginChanges(true)
	case 'r':
		ex.beginChanges(false)
	case 'g':
		ex.moveToTop()
	case 'G':
		ex.moveToBottom()
	case 'd':
		ex.beginDeletion()
	case 'j':
		ex.j()
	case 'k':
		ex.k()
	case ':':
		screen.SetFocus(0)
		ex.changeMode(mode.Command)
	case '/':
		ex.beginSearch()
	}

	if ex.action != searching {
		ex.fixCursor()
		ex.scroll()
	}
}

func (ex *Explorer) GetCursor(withBorder bool) (int, int, cursorType.CursorType) {
	var c cursorType.CursorType
	x := ex.visible.x + enums.InitialOffset
	y := ex.visible.y + enums.CursorOffset + 1

	switch ex.action {
	case changing:
		x += len(ex.entries[ex.cursor.y-searchBarOfset].name)
	case creating:
		x += len(ex.entries[len(ex.entries)-1].name)
		y += searchBarOfset
	case searching:
		x += len(ex.search)
		y = enums.CursorOffset
	}

	if !withBorder {
		x += 1
	}

	switch ex.action {
	case none:
		c = cursorType.CursorBloc
	case creating, changing, searching:
		c = cursorType.CursorLine
	case deleting:
		c = cursorType.CursorUnderline
	}

	return x, y, c
}

func (ex *Explorer) SetTitle() string {
	if slices.Equal(screen.Root(), ex.path) {
		return " Explorer "
	}
	return fmt.Sprintf(" %s ", filepath.Base(string(ex.path)))
}

func (ex *Explorer) Resize(w, h int) {
	ex.w = w
	ex.h = h
}

func (ex *Explorer) PreDraw() {
	cfg := config.GetConfig()
	if ex.showDot != cfg.ShowDot {
		ex.update = true
		ex.showDot = cfg.ShowDot
	}
	if ex.update || ex.action == searching {
		ex.scanEntries()
		ex.update = false
	}
	if ex.action == creating {
		ex.cursor.y = len(ex.entries) - 1
		ex.scroll()
	}
}
