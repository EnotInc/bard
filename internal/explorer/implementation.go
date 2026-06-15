package explorer

import (
	"fmt"
	"path/filepath"
	"slices"

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
	if ex.yScroll > 0 && index == 0 {
		return ascii.ArrowUp.Str()
	} else if index == ex.h-4 { // magic -4
		return ascii.ArrowDown.Str()
	}
	entry := ex.entries[index+ex.yScroll]
	var icon string
	if entry.isDir {
		icon = services.GetDirIcon(string(entry.name))
	} else {
		icon = services.GetFileIcon(string(entry.name))
	}
	e := fmt.Sprintf("%s%s", icon, string(ex.entries[index+ex.yScroll].name))
	e = services.VisibleSubString(e, 0, ex.w)
	return e
}

func (ex *Explorer) Handle(key rune) {
	if ex.typing {
		ex.typeNewEntry(key)
		return
	}

	switch key {
	case keys.Esc:
		screen.SetFocus(0)
	case keys.Enter:
		ex.openFileWithCallback()
		ex.cursor.y = 0
	case 'o':
		ex.typing = true
		ex.buffer = entry{name: []rune{}, isDir: false, path: ex.path}
		ex.cursor.y = len(ex.entries) - 1
	case 'r': // TODO: change file name (deletes it and let you type)
	case 'i': // TODO: change file name (set cursor to the end of the file name)
	case 'g':
		if slices.Equal(ex.root, ex.path) || len(ex.entries) == 0 {
			ex.cursor.y = 0
		} else {
			ex.cursor.y = 1
		}
	case 'G':
		ex.cursor.y = len(ex.entries) - 1
	case 'd':
		ex.delFileWithCallback()
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
	x := ex.visible.x + enums.InitialOffset + len(ex.buffer.name)
	y := ex.visible.y + enums.CursorOffset + 1

	if !withBorder {
		x += 1
	}

	return x, y
}

func (ex *Explorer) SetTitle() string {
	if slices.Equal(ex.root, ex.path) {
		return " Explorer "
	}
	return fmt.Sprintf(" %s ", filepath.Base(string(ex.path)))
}

func (ex *Explorer) Resize(w, h int) {
	ex.w = w
}

func (ex *Explorer) PreDraw() {
	ex.scanEntries()
	if ex.typing {
		ex.cursor.y = len(ex.entries) - 1
		ex.scroll()
	}
}
