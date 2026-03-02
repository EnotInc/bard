package editor

import "Enot/Bard/internal/mode"

func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b[e.curBuffer].InsertLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	case '\x7f':
		e.b[e.curBuffer].RemoveKey(0)
		e.ScrollLeft()
		e.ScrollUp()
	case '\t':
		//NOTE: yeah, I just insert 4 spaces instead of tabs
		for range 4 {
			e.b[e.curBuffer].InsertKey(' ')
			e.ScrollRight()
		}
	case '[', '{', '(', ')', '}', ']', '\'', '"', '<', '>', '*', '_', '`':
		e.b[e.curBuffer].InsertPair(key)
		e.ScrollRight()
	default:
		e.b[e.curBuffer].InsertKey(key)
		e.ScrollRight()
	}
	e.setUiCursor()
}
