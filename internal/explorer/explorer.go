package explorer

type Explorer struct {
	entries  []entry
	cursor   *cursor
	visible  *cursor
	w        int
	yScroll  int
	openFile func(file string)
	delFile  func(file string)
}

func InitExplorer(open func(file string), del func(file string), w int) *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
		w:        w,
		cursor:   c,
		visible:  v,
		openFile: open,
		delFile:  del,
	}
	ex.entries = scanEntries()
	ex.scroll()

	return ex
}
