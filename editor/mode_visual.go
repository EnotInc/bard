package editor

func (e *Editor) caseVisual(key rune) {
	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		e.subCmd += string(key)
	case 'h':
		e.moveWithSubCommand(e.b.H)
		e.setUiCursor()
		e.ScrollLeft()
	case 'j':
		e.moveWithSubCommand(e.b.J)
		e.ScrollDown()
		e.shiftLeft()
	case 'k':
		e.moveWithSubCommand(e.b.K)
		e.ScrollUp()
		e.shiftLeft()
	case 'l':
		e.moveWithSubCommand(e.b.L)
		e.ScrollRight()
	case 'y':
		e.b.copySelected(false, false)
		e.curMode = normal
	case 'x':
		e.b.copySelected(true, false)
		e.curMode = normal
	case 'd':
		e.b.copySelected(true, false)
		e.curMode = normal
	case '\033':
		e.curMode = normal
		e.ScrollLeft()
	}
}
