package editor

import (
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/mode"
)

// About caseVisualLine()
// Called from [Run()] func when current mode is Visual LIne
func (e *Editor) caseVisualLine(key rune) {
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
		e.b[e.curBuffer].CopySelected(false, true)
		e.curMode = mode.Normal
	case 'x':
		e.b[e.curBuffer].CopySelected(true, true)
		e.curMode = mode.Normal
	case 'd', 'D':
		e.b[e.curBuffer].CopySelected(true, true)
		e.curMode = mode.Normal
	case 's':
		e.b[e.curBuffer].CopySelected(true, true)
		e.b[e.curBuffer].InsertEmptyLine(enums.Above)
		e.b[e.curBuffer].MoveToFirstChar()
		e.curMode = mode.Insert
	case 'o', 'O':
		e.b[e.curBuffer].SwapTail()
	case 'w':
		e.b[e.curBuffer].MoveWord(1)
		e.setUiCursor()
	case 'W':
		e.b[e.curBuffer].MoveWORD(1)
		e.setUiCursor()
	case 'b':
		e.b[e.curBuffer].MoveBack(1)
		e.setUiCursor()
	case 'B':
		e.b[e.curBuffer].MoveBACK(1)
		e.setUiCursor()
	case 'e':
		e.b[e.curBuffer].MoveEnd(1)
		e.setUiCursor()
	case 'E':
		e.b[e.curBuffer].MoveEND(1)
		e.setUiCursor()
	case '\033':
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	}
}
