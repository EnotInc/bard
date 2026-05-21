package editor

import "strings"

func (e *Editor) IsGeneralMove(key rune) bool {
	// NOTE: this is not the best implementation I'm sure, but this is fine for not, ig
	return strings.Contains("webWEBhjklgG1234567890fFtT", string(key)) // && !(e.subCmd == "" && key == '0')
}

func (e *Editor) GeneralCase(key rune) {
	if e.subCmd == "" && key == '0' {
		e.b[e.curBuffer].MoveToFirstChar()
		return
	}
	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'f', 'F', 't', 'T':
		e.subCmd += string(key)
	case 'h':
		e.execWithSubCommand(e.b[e.curBuffer].H)
		e.setUiCursor()
		e.ScrollLeft()
	case 'j':
		e.execWithSubCommand(e.b[e.curBuffer].J)
		e.ScrollDown()
		e.shiftLeft()
	case 'k':
		e.execWithSubCommand(e.b[e.curBuffer].K)
		e.ScrollUp()
		e.shiftLeft()
	case 'l':
		e.execWithSubCommand(e.b[e.curBuffer].L)
		e.ScrollRight()
	case 'w':
		e.execWithSubCommand(e.b[e.curBuffer].MoveWord)
		e.setUiCursor()
	case 'W':
		e.execWithSubCommand(e.b[e.curBuffer].MoveWORD)
		e.setUiCursor()
	case 'b':
		e.execWithSubCommand(e.b[e.curBuffer].MoveBack)
		e.setUiCursor()
	case 'B':
		e.execWithSubCommand(e.b[e.curBuffer].MoveBACK)
		e.setUiCursor()
	case 'e':
		e.execWithSubCommand(e.b[e.curBuffer].MoveEnd)
		e.setUiCursor()
	case 'E':
		e.execWithSubCommand(e.b[e.curBuffer].MoveEND)
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
