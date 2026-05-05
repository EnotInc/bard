package editor

import "strings"

func IsGeneralMove(key rune) bool {
	// NOTE: this is not the best implementation I'm sure, but this is fine for not, ig
	return strings.Contains("webWEBhjklgG1234567890fFtT", string(key))
}

func (e *Editor) GeneralCase(key rune) {
	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'f', 'F', 't', 'T':
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
	case 'g':
		e.subCmd += "g"
		if e.subCmd == "gg" {
			e.b[e.curBuffer].MoveToFirstLine()
			e.setUiCursor()
			e.subCmd = ""
		}
	case 'G':
		e.b[e.curBuffer].MoveToLastLine()
		e.setUiCursor()
	}
}
