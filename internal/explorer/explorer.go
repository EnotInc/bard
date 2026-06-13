package explorer

import mode "github.com/EnotInc/Bard/internal/enums/mode"

type Explorer struct {
	cursor     *cursor
	visible    *cursor
	openFile   func(file string)
	delFile    func(file string)
	changeMode func(mode mode.Mode)
	path       string
	root       string
	buffer     entry
	entries    []entry
	w          int
	h          int
	yScroll    int
	typing     bool
}

const defaultRoot = "."
const back = ".."

func InitExplorer(open func(file string), del func(file string), change func(mode mode.Mode), w, h int) *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
		path:       defaultRoot,
		root:       defaultRoot,
		w:          w,
		h:          h,
		cursor:     c,
		visible:    v,
		openFile:   open,
		delFile:    del,
		changeMode: change,
		buffer:     entry{},
		typing:     false,
	}
	ex.scanEntries()
	ex.scroll()

	return ex
}

func (ex *Explorer) SetRoot(root string) {
	ex.path = root
	ex.root = root
}
