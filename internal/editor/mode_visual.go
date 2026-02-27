package editor

import "Enot/Bard/internal/mode"

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
		e.b.CopySelected(false, false)
		e.curMode = mode.Normal
	case 'x':
		e.b.CopySelected(true, false)
		e.curMode = mode.Normal
	case 'o', 'O':
		e.b.SwapTail()
	case 'd', 'D':
		e.b.CopySelected(true, false)
		e.curMode = mode.Normal
	case 's':
		e.b.CopySelected(true, false)
		e.curMode = mode.Insert
	case '\033':
		e.curMode = mode.Normal
		e.ScrollLeft()
	}
}
