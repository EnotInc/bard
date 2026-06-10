package explorer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"
)

type entry struct {
	name  string
	isDir bool
}

const depth = 0

func (ex *Explorer) scanEntries() {
	var e []entry
	var i = 0

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if i == 0 {
			i += 1
			return nil
		}
		if strings.Count(path, string(os.PathSeparator)) > depth {
			return fs.SkipDir
		}

		ent := entry{
			name:  d.Name(),
			isDir: d.IsDir(),
		}
		e = append(e, ent)

		i += 1
		return nil
	})

	if err != nil {
		panic("can't read current dir")
	}

	if ex.typing {
		e = append(e, ex.buffer)
	}
	ex.entries = e
}

func (ex *Explorer) openFileWithCallback() {
	entry := ex.entries[ex.cursor.y]
	if entry.isDir {
		return
	}
	ex.openFile(entry.name)
	screen.SendCall(calls.OpenFile)
}

func (ex *Explorer) delFileWithCallback() {
	entry := ex.entries[ex.cursor.y]
	ex.delFile(entry.name)
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
		}
	}
}

func (ex *Explorer) create(e entry) {
	if e.isDir {
		err := os.Mkdir(e.name, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		f, err := os.Create(e.name)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		ex.openFileWithCallback()
	}

}
