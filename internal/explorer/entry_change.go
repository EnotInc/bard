package explorer

import (
	"path/filepath"

	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/services"
)

func (ex *Explorer) beginChanges(clear bool) {
	ex.action = changing
	entry := ex.entries[ex.cursor.y]
	ex.buffer = entry
	entry.path = []rune(filepath.Join(string(ex.path), ""))

	if clear {
		ex.entries[ex.cursor.y].name = []rune{}
	}
}

func (ex *Explorer) handleChanges(key rune) {
	_entry := &ex.entries[ex.cursor.y]
	switch key {
	case keys.Esc:
		ex.entries[ex.cursor.y] = ex.buffer
		ex.buffer = entry{}
		ex.action = none
	case keys.Enter:
		ex.rename(string(ex.buffer.path), string(_entry.path))
		ex.action = none
		ex.update = true
	case keys.Backspace:
		if len(_entry.name) > 0 {
			_entry.name = _entry.name[:len(_entry.name)-1]
		}
	default:
		if services.IsLetterOrNumber(key) || key == '.' {
			_entry.name = append(_entry.name, key)
			_entry.path = []rune(filepath.Join(string(ex.path), string(_entry.name)))
		}
	}
}
