package explorer

import (
	"fmt"

	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/screen"
)

func (ex *Explorer) DrawLineAt(index int) string {
	if index+ex.yScroll >= len(ex.entries) {
		return "~"
	}
	return fmt.Sprintf("  %s", ex.entries[index+ex.yScroll].name)
}

func (ex *Explorer) Handle(key rune) {
	switch key {
	case keys.Esc:
		screen.SetFocus(0)
	case keys.Enter: // TODO: [un]fold dir
		ex.openFileWithCallback()
	case 'o': // TODO: create new file
	case 'r': // TODO: change file name (deletes it and let you type)
	case 'i': // TODO: change file name (set cursor to the end of the file name)
	case 'd':
		ex.delFileWithCallback()
	case 'j':
		ex.j()
	case 'k':
		ex.k()
	}
}

func (ex *Explorer) GetCursor(withBorder bool) (int, int) {
	borderOfset := 0
	if withBorder {
		borderOfset = 1
	}
	return ex.visible.x + enums.InitialOffset, ex.visible.y + enums.CursorOffset + borderOfset
}

func (ex *Explorer) SetTitle() string {
	return " Explorer "
}

func (ex *Explorer) Resize(w, h int) {}

func (ex *Explorer) PreDraw() {
	ex.entries = scanEntries()
}
