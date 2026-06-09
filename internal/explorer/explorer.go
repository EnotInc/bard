package explorer

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type entry struct {
	name  string
	isDir bool
}

type Explorer struct {
	entries []entry
	cursor  *cursor
	visible *cursor
	homeDir string
	yScroll int
}

func InitExplorer() *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
		homeDir: ".",
		cursor:  c,
		visible: v,
	}
	ex.entries = scanEntries()
	ex.scroll()

	return ex
}

const depth = 1

func scanEntries() []entry {
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

	return e
}
