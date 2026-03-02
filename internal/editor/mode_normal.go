package editor

import (
	"Enot/Bard/internal/mode"
	"strconv"
)

const (
	above = iota
	below
)

const (
	before = iota
	after
)

func (e *Editor) moveWithSubCommand(move func(int)) {
	if e.subCmd == "" {
		move(1)
		return
	}
	amount, err := strconv.Atoi(e.subCmd)
	if err != nil {
		e.subCmd = ""
		return
	}
	move(amount)
	e.subCmd = ""
}

func (e *Editor) replaceWithAmount(key rune) {
	if e.subCmd == "r" {
		e.caseReplaceChar(key, 1)
		return
	}

	amount, err := strconv.Atoi(e.subCmd[:len(e.subCmd)-1])
	if err != nil {
		e.subCmd = ""
		return
	}
	e.caseReplaceChar(key, amount)
	e.subCmd = ""
}

func (e *Editor) caseNormal(key rune) {
	cmd := []byte(e.subCmd)
	if len(cmd) > 0 && cmd[len(cmd)-1] == 'r' {
		e.replaceWithAmount(key)
		return
	}

	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'r':
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
	case 'i':
		e.curMode = mode.Insert
		e.ScrollLeft()
		e.tui.ShowHello = false
	case 'a':
		e.curMode = mode.Insert
		e.b.Insert_a()
		e.ScrollRight()
		e.tui.ShowHello = false
	case 'I':
		e.curMode = mode.Insert
		e.b.MoveToFirstVisible()
		e.moveLeft()
		e.tui.ShowHello = false
	case 'A':
		e.curMode = mode.Insert
		e.b.MoveToLastChar()
		e.moveRight()
		e.tui.ShowHello = false
	case ':':
		e.curMode = mode.Command
	case 'o':
		e.curMode = mode.Insert
		e.b.InsertEmptyLine(below)
		e.b.J(1)
		e.ScrollDown()
		e.moveLeft()
		e.tui.ShowHello = false
	case 'O':
		e.curMode = mode.Insert
		e.b.MoveToFirstChar()
		e.b.InsertEmptyLine(above)
		e.ScrollUp()
		e.moveLeft()
		e.tui.ShowHello = false
	case 'D':
		e.b.ClearLine()
		e.b.MoveToFirstChar()
		e.moveLeft()
	case 'd':
		e.subCmd += "d"
		if e.subCmd == "dd" {
			e.subCmd = ""
			e.b.RemoveLine()
			e.moveLeft()
		}
	case 'R':
		e.curMode = mode.Replace
	case 'x':
		e.b.Delkey()
		e.b.H(1)
		e.ScrollLeft()
	case 's':
		e.b.Delkey()
		e.curMode = mode.Insert
	case 'S':
		e.b.ClearLine()
		e.curMode = mode.Insert
	case 'g':
		e.subCmd += "g"
		if e.subCmd == "gg" {
			e.b.MoveToFirstLine()
			e.setUiCursor()
			e.subCmd = ""
		}
	case 'G':
		e.b.MoveToLastLine()
		e.setUiCursor()
	case 'w':
		e.b.MoveWord(1)
		e.setUiCursor()
	case 'W':
		e.b.MoveWORD(1)
		e.setUiCursor()
	case 'b':
		e.b.MoveBack(1)
		e.setUiCursor()
	case 'B':
		e.b.MoveBACK(1)
		e.setUiCursor()
	case 'e':
		e.b.MoveEnd(1)
		e.setUiCursor()
	case 'E':
		e.b.MoveEND(1)
		e.setUiCursor()
	case 'v':
		e.curMode = mode.Visual
		e.b.StartVisual()
		e.tui.ShowHello = false
	case 'V':
		e.curMode = mode.Visual_line
		e.b.StartVisualLine()
		e.tui.ShowHello = false
	case 'p':
		e.b.Paste(after)
		e.tui.ShowHello = false
	case 'P':
		e.b.Paste(before)
		e.tui.ShowHello = false
	default:
		e.subCmd = ""
	}
	e.setUiCursor()
}
