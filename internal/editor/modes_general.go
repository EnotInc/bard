package editor

import "strings"

func (e *Editor) IsGeneralMove(key rune) bool {
	// NOTE: this is not the best implementation I'm sure, but this is fine for not, ig
	return strings.Contains("webWEBhjklgG1234567890fFtT;", string(key)) &&
		!(len(e.subCmd) == 1 && strings.Contains("fFtT", e.subCmd))
}

func (e *Editor) replaceWith(key rune) bool {
	cmd := []rune(e.subCmd)
	if len(cmd) > 0 {
		if cmd[len(cmd)-1] == 'r' {
			e.replaceWithAmount(key)
			return true
		}
	}
	return false
}

func (e *Editor) findSome(key rune) bool {
	cmd := []byte(e.subCmd)
	if len(cmd) > 0 {
		switch cmd[len(cmd)-1] {
		case 'f':
			e.b[e.curBuffer].FindNext(key)
			e.lastCmd = e.subCmd + string(key)
			e.subCmd = ""
			return true

		case 'F':
			e.b[e.curBuffer].FindPrev(key)
			e.lastCmd = e.subCmd + string(key)
			e.subCmd = ""
			return true

		case 't':
			e.b[e.curBuffer].FindBeforeNext(key)
			e.lastCmd = e.subCmd + string(key)
			e.subCmd = ""
			return true

		case 'T':
			e.b[e.curBuffer].FindBeforePrev(key)
			e.lastCmd = e.subCmd + string(key)
			e.subCmd = ""
			return true
		}
	}
	return false
}

func (e *Editor) GeneralCase(key rune) {
	if e.subCmd == "" && key == '0' {
		e.b[e.curBuffer].MoveToFirstChar()
		return
	}
	if ok := e.replaceWith(key); ok {
		return
	}
	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'f', 'F', 't', 'T':
		e.subCmd += string(key)
	case ';':
		if len(e.lastCmd) <= 0 {
			return
		}

		e.subCmd = string(e.lastCmd[0])
		key := e.lastCmd[len(e.lastCmd)-1]
		if ok := e.findSome(rune(key)); !ok {
			e.tui.Message = "nothing was found"
		}
		return

	case 'h':
		e.execWithSubCommand(e.b[e.curBuffer].H)
		e.ScrollLeft()

	case 'j':
		e.execWithSubCommand(e.b[e.curBuffer].J)
		e.shiftLeft()

	case 'k':
		e.execWithSubCommand(e.b[e.curBuffer].K)
		e.shiftLeft()

	case 'l':
		e.execWithSubCommand(e.b[e.curBuffer].L)
		e.ScrollRight()

	case 'w':
		e.execWithSubCommand(e.b[e.curBuffer].MoveWord)

	case 'W':
		e.execWithSubCommand(e.b[e.curBuffer].MoveWORD)

	case 'b':
		e.execWithSubCommand(e.b[e.curBuffer].MoveBack)

	case 'B':
		e.execWithSubCommand(e.b[e.curBuffer].MoveBACK)

	case 'e':
		e.execWithSubCommand(e.b[e.curBuffer].MoveEnd)

	case 'E':
		e.execWithSubCommand(e.b[e.curBuffer].MoveEND)

	case 'g':
		e.subCmd += "g"
		if e.subCmd == "gg" {
			e.b[e.curBuffer].MoveToFirstLine()
			e.subCmd = ""
		}

	case 'G':
		e.b[e.curBuffer].MoveToLastLine()
	}
}
