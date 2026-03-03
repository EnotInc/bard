package editor

import "Enot/Bard/internal/mode"

func (e *Editor) caseVisual(key rune) {
	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		e.subCmd += string(key)
	case 'h':
		e.moveWithSubCommand(e.b[e.curBuffer].H)
		e.setUiCursor()
		e.ScrollLeft()
	case 'j':
		e.moveWithSubCommand(e.b[e.curBuffer].J)
		e.ScrollDown()
		e.shiftLeft()
	case 'k':
		e.moveWithSubCommand(e.b[e.curBuffer].K)
		e.ScrollUp()
		e.shiftLeft()
	case 'l':
		e.moveWithSubCommand(e.b[e.curBuffer].L)
		e.ScrollRight()
	case 'y':
		e.b[e.curBuffer].CopySelected(false, false)
		e.curMode = mode.Normal
	case 'x':
		e.b[e.curBuffer].CopySelected(true, false)
		e.curMode = mode.Normal
	case 'o', 'O':
		e.b[e.curBuffer].SwapTail()
	case 'd', 'D':
		e.b[e.curBuffer].CopySelected(true, false)
		e.curMode = mode.Normal
	case 's':
		e.b[e.curBuffer].CopySelected(true, false)
		e.curMode = mode.Insert
	case '\033':
		e.curMode = mode.Normal
		e.ScrollLeft()
	}
}
