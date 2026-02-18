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

func (e *Editor) caseNormal(key rune) {
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
	case 'i':
		e.curMode = insert
		e.ScrollLeft()
		e.showInfo = false
	case 'a':
		e.curMode = insert
		if len(e.b.lines[e.b.cursor.line].data) > 0 {
			e.b.cursor.offset += 1
		}
		e.ScrollRight()
		e.showInfo = false
	case 'I':
		e.curMode = insert
		e.b.moveToFirst()
		e.moveLeft()
		e.showInfo = false
	case 'A':
		e.curMode = insert
		e.b.cursor.offset = len(e.b.lines[e.b.cursor.line].data)
		e.moveRight()
		e.showInfo = false
	case ':':
		e.curMode = command
	case 'o':
		e.curMode = insert
		e.b.cursor.offset = 0
		e.b.InsertEmptyLine(below)
		e.b.cursor.line += 1
		e.ScrollDown()
		e.moveLeft()
		e.showInfo = false
	case 'O':
		e.curMode = insert
		e.b.cursor.offset = 0
		e.b.InsertEmptyLine(above)
		e.ScrollUp()
		e.moveLeft()
		e.showInfo = false
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
	case 'x':
		e.b.Delkey()
		if e.b.cursor.offset >= len(e.b.lines[e.b.cursor.line].data) && e.b.cursor.offset > 0 {
			e.b.cursor.offset -= 1
		}
		e.ScrollLeft()
	case 's':
		e.b.Delkey()
		e.curMode = insert
	case 'v':
		e.curMode = visual
		e.b.visual.line = e.b.cursor.line
		e.b.visual.offset = e.b.cursor.offset
		e.showInfo = false
	case 'V':
		e.curMode = visual_line
		e.b.visual.line = e.b.cursor.line
		e.showInfo = false
	case 'p':
		e.b.paste(after)
		e.showInfo = false
	case 'P':
		e.b.paste(before)
		e.showInfo = false
	default:
		e.subCmd = ""
	}
	e.setUiCursor()
}
