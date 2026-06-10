package explorer

type Explorer struct {
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

func InitExplorer(open func(file string), del func(file string), w, h int) *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
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
