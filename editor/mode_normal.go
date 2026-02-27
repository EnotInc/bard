package editor

import (
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
		e.curMode = insert
		e.ScrollLeft()
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'a':
		e.curMode = insert
		if len(e.b.lines[e.b.cursor.line].data) > 0 {
			e.b.cursor.offset += 1
		}
		e.ScrollRight()
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'I':
		e.curMode = insert
		e.b.moveToFirst()
		e.moveLeft()
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'A':
		e.curMode = insert
		e.b.cursor.offset = len(e.b.lines[e.b.cursor.line].data)
		e.moveRight()
		e.showHello = false
		e.ui.hello = [][]rune{}
	case ':':
		e.curMode = command
	case 'o':
		e.curMode = insert
		e.b.cursor.offset = 0
		e.b.InsertEmptyLine(below)
		e.b.cursor.line += 1
		e.ScrollDown()
		e.moveLeft()
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'O':
		e.curMode = insert
		e.b.cursor.offset = 0
		e.b.InsertEmptyLine(above)
		e.ScrollUp()
		e.moveLeft()
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'D':
		e.b.ClearLine()
		e.b.cursor.offset = 0
		e.moveLeft()
	case 'd':
		e.subCmd += "d"
		if e.subCmd == "dd" {
			e.subCmd = ""
			e.b.RemoveLine()
			e.moveLeft()
		}
	case 'R':
		e.curMode = replace
	case 'x':
		e.b.Delkey()
		if e.b.cursor.offset >= len(e.b.lines[e.b.cursor.line].data) && e.b.cursor.offset > 0 {
			e.b.cursor.offset -= 1
		}
		e.ScrollLeft()
	case 's':
		e.b.Delkey()
		e.curMode = insert
	case 'S':
		e.b.ClearLine()
		e.curMode = insert
	case 'g':
		e.subCmd += "g"
		if e.subCmd == "gg" {
			e.b.moveToFirstLine()
			e.setUiCursor()
			e.subCmd = ""
		}
	case 'G':
		e.b.moveToLastLine()
		e.setUiCursor()
	case 'w':
		e.b.moveWord(1)
		e.setUiCursor()
	case 'W':
		e.b.moveWORD(1)
		e.setUiCursor()
	case 'b':
		e.b.moveBack(1)
		e.setUiCursor()
	case 'e':
		e.b.moveEnd(1)
		e.setUiCursor()
	case 'v':
		e.curMode = visual
		e.b.visual.line = e.b.cursor.line
		e.b.visual.offset = e.b.cursor.offset
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'V':
		e.curMode = visual_line
		e.b.visual.line = e.b.cursor.line
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'p':
		e.b.paste(after)
		e.showHello = false
		e.ui.hello = [][]rune{}
	case 'P':
		e.b.paste(before)
		e.showHello = false
		e.ui.hello = [][]rune{}
	default:
		e.subCmd = ""
	}
	e.setUiCursor()
}
