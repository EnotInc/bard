package explorer

import (
	"os"
	"path/filepath"

	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/services"
)

func (ex *Explorer) beginCreation() {
	ex.action = creating

	entry := entry{name: []rune{}, path: ex.path, isDir: false}
	ex.entries = append(ex.entries, entry)
	ex.cursor.y = len(ex.entries) + searchBarOfset
}

func (ex *Explorer) handleCreate(key rune) {
	_entry := &ex.entries[len(ex.entries)-1]
	switch key {
	case keys.Esc:
		ex.entries = ex.entries[:len(ex.entries)-1]
		ex.moveToBottom()
		ex.scroll()
		ex.action = none
		ex.update = true
	case keys.Enter:
		ex.create(_entry)
		ex.moveToTop()
		ex.scroll()
		ex.action = none
		ex.update = true
	case keys.Backspace:
		if len(_entry.name) > 0 {
			_entry.name = _entry.name[:len(_entry.name)-1]
		}
	default:
		if key == '/' {
			_entry.isDir = !_entry.isDir
			return
		}
		if services.IsLetterOrNumber(key) || key == '.' {
			_entry.name = append(_entry.name, key)
			_entry.path = []rune(filepath.Join(string(ex.path), string(_entry.name)))
		}
	}
}

func (ex *Explorer) create(e *entry) {
	if e.isDir {
		err := os.Mkdir(string(e.path), 0755)
		if err != nil {
			ex.setError(err.Error())
		}
	} else {
		f, err := os.Create(string(e.path))
		if err != nil {
			ex.setError(err.Error())
		}
		defer f.Close()
	}
	ex.openEntry(e)
}
