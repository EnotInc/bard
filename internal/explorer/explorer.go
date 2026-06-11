package explorer

type Explorer struct {
	root     string
	entries  []entry
	cursor   *cursor
	visible  *cursor
	w, h     int
	yScroll  int
	openFile func(file string)
	delFile  func(file string)
	buffer   entry
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
