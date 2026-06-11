package explorer

type Explorer struct {
	cursor   *cursor
	visible  *cursor
	openFile func(file string)
	delFile  func(file string)
	root     string
	buffer   entry
	entries  []entry
	w        int
	h        int
	yScroll  int
	typing   bool
}

const root = "."
const back = ".."

func InitExplorer(open func(file string), del func(file string), w, h int) *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
		root:     root,
		w:        w,
		h:        h,
		cursor:   c,
		visible:  v,
		openFile: open,
		delFile:  del,
		buffer:   entry{},
		typing:   false,
	}
	ex.scanEntries()
	ex.scroll()

	return ex
}
