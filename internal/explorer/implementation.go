package explorer

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"
)

func (ex *Explorer) DrawLineAt(index int) string {
	if index+ex.yScroll >= len(ex.entries) {
		return ""
	}

	cfg := config.GetConfig()
	ofset := 2
	if cfg.ShowBorder {
		ofset *= 2
	}

	if ex.yScroll > 0 && index == 0 {
		return ascii.ArrowUp.Str()
	} else if index == ex.h-ofset {
		return ascii.ArrowDown.Str()
	}

	entry := ex.entries[index+ex.yScroll]
	var icon string
	if entry.isDir {
		icon = services.GetDirIcon(string(entry.name))
	} else {
		icon = services.GetFileIcon(string(entry.name))
	}
	if ex.action == deleting && index == ex.visible.y+1 {
		icon = fmt.Sprintf(" y/n: %s%s", ascii.Stricked, icon)
	}

	e := fmt.Sprintf("%s%s", icon, string(ex.entries[index+ex.yScroll].name))
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
	}

	switch key {
	case keys.Esc:
		screen.SetFocus(0)
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
	}
	ex.fixCursor()
	ex.scroll()
}

func (ex *Explorer) GetCursor(withBorder bool) (int, int) {
	x := ex.visible.x + enums.InitialOffset
	y := ex.visible.y + enums.CursorOffset + 1

	if ex.action != none && ex.action != deleting {
		x += len(ex.entries[ex.cursor.y].name)
	}

	if !withBorder {
		x += 1
	}

	return x, y
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
	if ex.update {
		ex.scanEntries()
		ex.update = false
	}
	if ex.action == creating {
		ex.cursor.y = len(ex.entries) - 1
		ex.scroll()
	}
}
