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
		e.tui.Hello = [][]rune{}
	case 'a':
		e.curMode = mode.Insert
		if len(e.b.Lines[e.b.Cursor.Line].Data) > 0 {
			e.b.Cursor.Offset += 1
		}
		e.ScrollRight()
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case 'I':
		e.curMode = mode.Insert
		e.b.MoveToFirst()
		e.moveLeft()
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case 'A':
		e.curMode = mode.Insert
		e.b.Cursor.Offset = len(e.b.Lines[e.b.Cursor.Line].Data)
		e.moveRight()
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case ':':
		e.curMode = mode.Command
	case 'o':
		e.curMode = mode.Insert
		e.b.Cursor.Offset = 0
		e.b.InsertEmptyLine(below)
		e.b.Cursor.Line += 1
		e.ScrollDown()
		e.moveLeft()
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case 'O':
		e.curMode = mode.Insert
		e.b.Cursor.Offset = 0
		e.b.InsertEmptyLine(above)
		e.ScrollUp()
		e.moveLeft()
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case 'D':
		e.b.ClearLine()
		e.b.Cursor.Offset = 0
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
	case 'e':
		e.b.MoveEnd(1)
		e.setUiCursor()
	case 'v':
		e.curMode = mode.Visual
		e.b.Visual.Line = e.b.Cursor.Line
		e.b.Visual.Offset = e.b.Cursor.Offset
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case 'V':
		e.curMode = mode.Visual_line
		e.b.Visual.Line = e.b.Cursor.Line
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case 'p':
		e.b.Paste(after)
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	case 'P':
		e.b.Paste(before)
		e.tui.ShowHello = false
		e.tui.Hello = [][]rune{}
	default:
		e.subCmd = ""
	}
	e.setUiCursor()
}
