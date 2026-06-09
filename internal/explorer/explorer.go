package explorer

type Explorer struct {
	entries []entry
	cursor  *cursor
	visible *cursor
	yScroll int
}

func InitExplorer() *Explorer {

	c := initCursor()
	v := initCursor()
	ex := &Explorer{
		cursor:  c,
		visible: v,
	}
	ex.entries = scanEntries()
	ex.scroll()

	return ex
}
