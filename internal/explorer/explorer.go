package explorer

type Explorer struct {
	entries  []entry
	cursor   *cursor
	visible  *cursor
	yScroll  int
	openFile func(file string)
	delFile  func(file string)
}

func InitExplorer(open func(file string), del func(file string)) *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
		cursor:   c,
		visible:  v,
		openFile: open,
		delFile:  del,
	}
	ex.entries = scanEntries()
	ex.scroll()

	return ex
}
