package editor

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	initialLineShift = 1
	ScrollBorder     = 5
)

const (
	above = iota
	below
)

const (
	clearView    = "\033[2J"
	clearHistory = "\033[3J"
	moveToStart  = "\033[0H"
	cursorYellow = "\033]12;yellow\x07"
	cursorReset  = "\033]112\a"
)

type Mode string

const (
	normal      Mode = "NORMAL"
	command     Mode = "COMMAND"
	insert      Mode = "INSERT"
	visual      Mode = "VISUAL"
	visual_line Mode = "VISUAL-LINE"
)

type Editor struct {
	oldState   *term.State
	b          *Buffer
	ui         *UI
	curMode    Mode
	curCommand string
	file       string
	fdIn       int
	w          int
	h          int
}

func InitEditor() *Editor {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil {
		panic(err)
	}
	_w, _h, _ := term.GetSize(_fdOut)

	_b := InitBuffer()
	_ui := InitUI(_h, _w)

	e := &Editor{
		oldState:   old,
		b:          _b,
		ui:         _ui,
		curMode:    normal,
		curCommand: "",
		fdIn:       _fdIn,
		w:          _w,
		h:          _h,
	}

	fmt.Print(cursorYellow)
	return e
}

func (e *Editor) ScrollUp() {
	if e.ui.curRow == ScrollBorder {
		if e.ui.upperBorder != 0 {
			e.ui.upperBorder -= 1
			e.ui.lowerBorder -= 1
		}
		e.setUiCursor()
	}
}

func (e *Editor) ScrollDown() {
	if e.ui.curRow == e.h-ScrollBorder {
		if e.ui.lowerBorder != len(e.b.lines)+ScrollBorder {
			e.ui.upperBorder += 1
			e.ui.lowerBorder += 1
		}
		e.setUiCursor()
	}
}

func (e *Editor) ScrollRight() {
	if e.ui.curOff >= e.w-ScrollBorder*2 {
		if e.ui.rightBorder != len(e.b.lines[e.b.cursor.line].data)+ScrollBorder {
			e.ui.leftBorder += 1
			e.ui.rightBorder += 1
		}
		e.setUiCursor()
	}
}

func (e *Editor) ScrollLeft() {
	if e.ui.curOff <= ScrollBorder {
		if e.ui.leftBorder != 0 {
			e.ui.leftBorder -= 1
			e.ui.rightBorder -= 1
		}
		e.setUiCursor()
	}
}

func (e *Editor) setUiCursor() {
	e.ui.curRow = e.b.cursor.line - e.ui.upperBorder
	e.ui.curOff = e.b.cursor.ofset - e.ui.leftBorder
}

func (e *Editor) moveLeft() {
	e.ui.leftBorder = 0
	e.ui.rightBorder = e.w - initialOfset
	//e.ui.curOff = e.b.cursor.ofset - e.ui.leftBorder
	e.setUiCursor()
}

func (e *Editor) moveRight() {
	e.ui.rightBorder = len(e.b.lines[e.b.cursor.line].data) + ScrollBorder*2
	e.ui.leftBorder = e.ui.rightBorder - e.w
	if e.ui.leftBorder < 0 {
		e.ui.leftBorder = 0
	}
	if e.ui.rightBorder < e.w-initialOfset {
		e.ui.rightBorder = e.w - initialOfset
	}
	e.setUiCursor()
}

func (e *Editor) caseNormal(key rune) {
	//e.b.cursor.lastOfset = e.b.cursor.ofset
	switch key {
	case 'h':
		e.b.H()
		e.ScrollLeft()
	case 'j':
		e.b.J()
		e.ScrollDown()
	case 'k':
		e.b.K()
		e.ScrollUp()
	case 'l':
		e.b.L()
		e.ScrollRight()
	case 'v':
		e.curMode = visual
	case 'i':
		e.curMode = insert
		e.ScrollLeft()
	case 'a':
		e.curMode = insert
		if len(e.b.lines[e.b.cursor.line].data) > 0 {
			e.b.cursor.ofset += 1
		}
		e.ScrollRight()
	case 'I':
		e.curMode = insert
		for i := range len(e.b.lines[e.b.cursor.line].data) {
			if e.b.lines[e.b.cursor.line].data[i] != ' ' {
				e.b.cursor.ofset = i
				break
			}
		}
		e.moveLeft()
	case 'A':
		e.curMode = insert
		e.b.cursor.ofset = len(e.b.lines[e.b.cursor.line].data)
		e.moveRight()
	case ':':
		e.curMode = command
	case 'o':
		e.curMode = insert
		e.b.cursor.ofset = 0
		e.b.InsertEmptyLine(below)
		e.b.cursor.line += 1
		e.ScrollDown()
		e.moveLeft()
	case 'O':
		e.curMode = insert
		e.b.cursor.ofset = 0
		e.b.InsertEmptyLine(above)
		e.ScrollUp()
		e.moveLeft()
	case 'x':
		//e.b.Yank()
		e.b.Delkey()
		if e.b.cursor.ofset >= len(e.b.lines[e.b.cursor.line].data) && e.b.cursor.ofset > 0 {
			e.b.cursor.ofset -= 1
		}
		e.ScrollLeft()
	//TODO: reimplement
	// case 'p':
	// 	e.b.Paste()
	case 's':
		e.b.Delkey()
		e.curMode = insert
	}
	e.setUiCursor()
}

func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b.InsertLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = normal
		if e.b.cursor.ofset > 0 {
			e.b.cursor.ofset -= 1
		}
		e.ScrollLeft()
	case '\127', '\x7f':
		e.b.RemoveKey(0)
		e.ScrollLeft()
		e.ScrollUp()
	case '\t':
		//NOTE: yeah, I just insert 4 spaces instead of tabs
		for range 4 {
			e.b.InsertKey(' ')
			e.ScrollRight()
		}
	default:
		e.b.InsertKey(key)
		e.ScrollRight()
	}
	e.setUiCursor()
}

func (e *Editor) caseCommand(key rune) {
	switch key {
	case '\033':
		e.curCommand = ""
		e.curMode = normal
	case '\127', '\x7f':
		if len(e.curCommand) > 0 {
			e.curCommand = e.curCommand[:len(e.curCommand)-1]
		} else {
			e.curCommand = ""
			e.curMode = normal
		}
	case '\013', '\r', '\n':
		success := e.execCommand()
		if !success {
			e.curCommand = ""
			e.curMode = normal
		}
	default:
		e.curCommand += string(key)
	}
}

func (e *Editor) execCommand() bool {
	switch e.curCommand {
	case "q":
		fmt.Print(clearView, clearHistory, moveToStart, cursorReset)
		term.Restore(e.fdIn, e.oldState)
		os.Exit(0)
		return true
	case "w":
		//TODO: add notification if there is no file name provided
		err := e.SaveFile()
		if err != nil {
			return false
		}
		e.curCommand = ""
		e.curMode = normal
		return true
	case "x", "wq":
		err := e.SaveFile()
		if err != nil {
			return false
		} else {
			fmt.Print(clearView, clearHistory, moveToStart, cursorReset)
			term.Restore(e.fdIn, e.oldState)
			os.Exit(0)
			return true
		}
	case "rln":
		e.ui.rln = !e.ui.rln
		e.curCommand = ""
		e.curMode = normal
		return true
	default:
		return false
	}
}

func (e *Editor) caseVisual(key rune) {
	switch key {
	case '\033':
		e.curMode = normal
	case 'h':
		e.b.H()
	case 'j':
		e.b.J()
		e.ScrollDown()
	case 'k':
		e.b.K()
		e.ScrollUp()
	case 'l':
		e.b.L()
	}
}

func (e *Editor) Run() {
	e.ui.Draw(e)
	reader := bufio.NewReader(os.Stdin)
	for {
		key, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}
		switch e.curMode {
		case normal:
			e.caseNormal(key)
		case command:
			e.caseCommand(key)
		case insert:
			e.caseInsert(key)
		case visual:
			e.caseVisual(key)
		}

		e.ui.Draw(e)
	}
}
