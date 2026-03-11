package editor

import (
	"Enot/Bard/internal/enums"
	"Enot/Bard/internal/mode"
)

func (e *Editor) caseNormal(key rune) {
	cmd := []byte(e.subCmd)
	if len(cmd) > 0 {
		switch cmd[len(cmd)-1] {
		case 'r':
			e.replaceWithAmount(key)
			return
		case 'f':
			e.b[e.curBuffer].FindNext(key)
			e.subCmd = ""
			return
		case 'F':
			e.b[e.curBuffer].FindPrev(key)
			e.subCmd = ""
			return
		case 't':
			e.b[e.curBuffer].FindBeforeNext(key)
			e.subCmd = ""
			return
		case 'T':
			e.b[e.curBuffer].FindBeforePrev(key)
			e.subCmd = ""
			return
		}
	}

	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'r', 'f', 'F', 't', 'T':
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
	case 'i':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
		e.ScrollLeft()
		e.tui.ShowHello = false
	case 'a':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
		e.b[e.curBuffer].FixOffset()
		e.b[e.curBuffer].Insert_a()
		e.ScrollRight()
		e.tui.ShowHello = false
	case 'I':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
		e.b[e.curBuffer].MoveToFirstVisible()
		e.moveLeft()
		e.tui.ShowHello = false
	case 'A':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
		e.b[e.curBuffer].MoveToLastChar()
		e.moveRight()
		e.tui.ShowHello = false
	case ':':
		e.curMode = mode.Command
	case 'o':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
		e.b[e.curBuffer].InsertEmptyLine(enums.Below)
		e.b[e.curBuffer].J(1)
		e.ScrollDown()
		e.moveLeft()
		e.tui.ShowHello = false
	case 'O':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
		e.b[e.curBuffer].InsertEmptyLine(enums.Above)
		e.b[e.curBuffer].MoveToFirstChar()
		e.ScrollUp()
		e.moveLeft()
		e.tui.ShowHello = false
	case 'D':
		e.b[e.curBuffer].ClearLine()
		e.b[e.curBuffer].MoveToFirstChar()
		e.moveLeft()
	case 'd':
		e.subCmd += "d"
		if e.subCmd == "dd" {
			e.subCmd = ""
			e.b[e.curBuffer].RemoveLine()
			e.moveLeft()
		}
	case 'R':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Replace
		}
	case 'x':
		e.b[e.curBuffer].Delkey()
		e.b[e.curBuffer].H(1)
		e.ScrollLeft()
	case 's':
		e.b[e.curBuffer].Delkey()
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
	case 'S':
		e.b[e.curBuffer].ClearLine()
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
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
	case 'v':
		e.curMode = mode.Visual
		e.b[e.curBuffer].StartVisual()
		e.tui.ShowHello = false
	case 'V':
		e.curMode = mode.Visual_line
		e.b[e.curBuffer].StartVisualLine()
		e.tui.ShowHello = false
	case 'p':
		e.b[e.curBuffer].Paste(enums.After)
		e.tui.ShowHello = false
	case 'P':
		e.b[e.curBuffer].Paste(enums.Before)
		e.tui.ShowHello = false
	default:
		e.subCmd = ""
	}
	e.setUiCursor()
}
