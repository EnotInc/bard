package explorer

type Explorer struct {
	entries  []entry
	cursor   *cursor
	visible  *cursor
	yScroll  int
	openFile func(file string)
}

func InitExplorer(callback func(file string)) *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
		cursor:   c,
		visible:  v,
		openFile: callback,
	}
	ex.entries = scanEntries()
	ex.scroll()

	return ex
}
