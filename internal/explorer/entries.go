package explorer

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"
)

type entry struct {
	name  []rune
	path  []rune
	isDir bool
}

func (ex *Explorer) scanEntries() {
	cfg := config.GetConfig()
	var e []entry
	entries, err := os.ReadDir(string(ex.path))
	if err != nil {
		panic(err)
	}

	if len(entries) == 0 && !ex.typing { // cur dir is added if dir is empty
		e = append(e, entry{name: []rune(defaultRoot), isDir: true})
	}

	if !slices.Equal(ex.root, ex.path) { // 'go back' entry
		e = append(e, entry{name: []rune(back), isDir: true})
	}

	for _, en := range entries {
		if strings.HasPrefix(en.Name(), ".") && !cfg.ShowDot {
			continue
		}
		ent := entry{
			name:  []rune(en.Name()),
			path:  []rune(filepath.Join(string(ex.path), en.Name())),
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
		if slices.Equal(entry.name, []rune(back)) {
			ex.path = []rune(filepath.Dir(string(ex.path)))
		} else {
			ex.path = []rune(filepath.Join(string(ex.path), string(entry.name)))
		}
		return
	}
	ex.openFile(string(entry.path))
	screen.SendCall(calls.OpenFile)
}

func (ex *Explorer) delFileWithCallback() {
	entry := ex.entries[ex.cursor.y]
	ex.delFile(string(entry.path))
	screen.SendCall(calls.DelFile)
}

func (ex *Explorer) typeNewEntry(key rune) {
	switch key {
	case keys.Esc:
		ex.typing = false
		ex.buffer = entry{}
		ex.entries = ex.entries[:len(ex.entries)-1]
		ex.cursor.y = max(len(ex.entries)-1, 0)
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
			ex.buffer.name = append(ex.buffer.name, key)
			ex.buffer.path = []rune(filepath.Join(string(ex.path), string(ex.buffer.name)))
		}
	}
}

func (ex *Explorer) create(e entry) {
	if e.isDir {
		err := os.Mkdir(string(e.path), 0755)
		if err != nil {
			panic(err)
		}
	} else {
		f, err := os.Create(string(e.path))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		ex.openFileWithCallback()
	}

}
