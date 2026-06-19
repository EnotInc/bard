package explorer

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/screen"
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

	if !slices.Equal(screen.Root(), ex.path) {
		e = append(e, entry{name: []rune(enums.Back), isDir: true})
	}

	if len(e) == 0 && len(entries) == 0 {
		e = append(e, entry{name: []rune(enums.DefaultRoot), isDir: true})
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

	ex.entries = e
}

func (ex *Explorer) openEntryWithCallback() {
	entry := ex.entries[ex.cursor.y]
	if entry.isDir {
		if slices.Equal(entry.name, []rune(enums.Back)) {
			ex.path = []rune(filepath.Dir(string(ex.path)))
		} else if !slices.Equal(entry.name, []rune(enums.DefaultRoot)) {
			ex.path = []rune(filepath.Join(string(ex.path), string(entry.name)))
		}
		ex.update = true
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
