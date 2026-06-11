package explorer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"
)

type entry struct {
	name  string
	path  string
	isDir bool
}

func (ex *Explorer) scanEntries() {
	var e []entry
	entries, err := os.ReadDir(ex.root)
	if err != nil {
		panic(err)
	}
	if ex.curPath != ex.root {
		e = append(e, entry{name: back, isDir: true})
	}
	for _, en := range entries {
		ent := entry{
			name:  en.Name(),
			path:  filepath.Join(ex.root, en.Name()),
			isDir: en.IsDir(),
		}
		e = append(e, ent)
	}

	if ex.typing {
		e = append(e, ex.buffer)
	}
	ex.entries = e
}

func (ex *Explorer) openFileWithCallback() {
	entry := ex.entries[ex.cursor.y]
	if entry.isDir {
		if entry.name == back {
			ex.root = filepath.Dir(ex.root)
		} else {
			ex.root = filepath.Join(ex.root, entry.name)
		}
		return
	}
	ex.openFile(entry.path)
	screen.SendCall(calls.OpenFile)
}

func (ex *Explorer) delFileWithCallback() {
	entry := ex.entries[ex.cursor.y]
	ex.delFile(entry.path)
	screen.SendCall(calls.DelFile)
}

func (ex *Explorer) typeNewEntry(key rune) {
	switch key {
	case keys.Esc:
		ex.typing = false
		ex.buffer = entry{}
		ex.entries = ex.entries[:len(ex.entries)-1]
		ex.cursor.y = len(ex.entries) - 1
		ex.scroll()
	case keys.Enter:
		ex.typing = false
		ex.create(ex.buffer)
		ex.buffer = entry{}
		ex.entries = ex.entries[:len(ex.entries)-1]
		ex.cursor.y = 0
		ex.scroll()
	case keys.Backspace:
		if len(ex.buffer.name) > 0 {
			ex.buffer.name = ex.buffer.name[:len(ex.buffer.name)-1]
		}
	default:
		if key == '/' {
			ex.buffer.isDir = !ex.buffer.isDir
			return
		}
		if services.IsLetterOrNumber(key) || key == '.' {
			ex.buffer.name = fmt.Sprintf("%s%c", ex.buffer.name, key)
			ex.buffer.path = filepath.Join(ex.root, ex.buffer.name)
		}
	}
}

func (ex *Explorer) create(e entry) {
	if e.isDir {
		err := os.Mkdir(e.path, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		f, err := os.Create(e.path)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		ex.openFileWithCallback()
	}

}
