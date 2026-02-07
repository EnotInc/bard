package editor

func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b.InsertLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = normal
		if e.b.cursor.ofset > 0 {
			e.b.cursor.ofset -= 1
		}
		e.ScrollLeft()
	case '\x7f':
		e.b.RemoveKey(0)
		e.ScrollLeft()
		e.ScrollUp()
	case '\t':
		//NOTE: yeah, I just insert 4 spaces instead of tabs
		for range 4 {
			e.b.InsertKey(' ')
			e.ScrollRight()
		}
	default:
		e.b.InsertKey(key)
		e.ScrollRight()
	}
	e.setUiCursor()
}
