package editor

import "Enot/Bard/internal/mode"

func (e *Editor) caseReplaceChar(key rune, amount int) {
	switch key {
	case '\013', '\r', '\n':
		e.b.InsertLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = mode.Normal
		e.b.EscapeToNormal()
		e.ScrollLeft()
	case '\x7f': // just do nothing if backspace is pressed
		return
	case '\t':
		e.b.Delkey()
		for range 4 {
			e.b.InsertKey(' ')
			e.ScrollRight()
		}
	default:
		e.b.ReplaceKeys(key, amount)
	}
	e.subCmd = ""
}

func (e *Editor) caseReplaceMode(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b.DelAndMoveLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = mode.Normal
		e.b.EscapeToNormal()
		e.ScrollLeft()
	case '\x7f':
		e.b.RemoveKey(0)
		e.ScrollLeft()
		e.ScrollUp()
	case '\t':
		e.b.Delkey()
		for range 4 {
			e.b.InsertKey(' ')
			e.ScrollRight()
		}
	default:
		e.b.ReplaceKeys(key, 1)
		e.b.L(1)
	}
}
