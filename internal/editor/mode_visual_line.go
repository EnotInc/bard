package editor

import "Enot/Bard/internal/mode"

func (e *Editor) caseVisualLine(key rune) {
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
		e.b.CopySelected(false, true)
		e.curMode = mode.Normal
	case 'x':
		e.b.CopySelected(true, true)
		e.curMode = mode.Normal
	case 'd', 'D':
		e.b.CopySelected(true, true)
		e.curMode = mode.Normal
	case 's':
		e.b.CopySelected(true, true)
		e.b.InsertEmptyLine(above)
		e.b.MoveToFirst()
		e.curMode = mode.Insert
	case 'o', 'O':
		e.b.SwapTail()
	case '\033':
		e.curMode = mode.Normal
		if e.b.Cursor.Offset > 0 {
			e.b.Cursor.Offset -= 1
		}
		e.ScrollLeft()
	}
}
