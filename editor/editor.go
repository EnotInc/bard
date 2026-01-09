package editor

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	initialCurOfset  = 2
	initialLineShift = 1
)

const (
	above = iota
	below
)

const (
	clearView    = "\033[2J"
	clearHistory = "\033[3J"
	moveToStart  = "\033[0H"
	cursorBloc   = "\x1b[0 q"
	cursorLine   = "\x1b[5 q"
)

type Mode string

const (
	normal  = "NORMAL"
	command = "COMMAND"
	insert  = "INSERT"
)

type Editor struct {
	oldState   *term.State
	b          *Buffer
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
	w, h, _ := term.GetSize(_fdOut)

	b := InitBuffer()

	e := &Editor{
		oldState:   old,
		b:          b,
		curMode:    normal,
		curCommand: "",
		fdIn:       _fdIn,
		w:          w,
		h:          h,
	}

	return e
}

func (e *Editor) caseNormal(key rune) {
	switch key {
	case 'h':
		e.b.H()
	case 'j':
		e.b.J()
	case 'k':
		e.b.K()
	case 'l':
		e.b.L()
	case 'i':
		e.curMode = insert
	case 'a':
		e.curMode = insert
		if len(e.b.lines[e.b.cursor.line].data) > 0 {
			e.b.cursor.ofset += 1
		}
	case 'I':
		e.curMode = insert
		//TODO: move to the first char instad of the 0
		e.b.cursor.ofset = 0
	case 'A':
		e.curMode = insert
		e.b.cursor.ofset = len(e.b.lines[e.b.cursor.line].data)
	case ':':
		e.curMode = command
	case 'o':
		e.curMode = insert
		e.b.cursor.ofset = 0
		e.b.InsertEmptyLine(below)
		e.b.cursor.line += 1
	case 'O':
		e.curMode = insert
		e.b.cursor.ofset = 0
		e.b.InsertEmptyLine(above)
	case 'x':
		e.b.Delkey()
	}
}

func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b.InsertLine()
	case '\033':
		e.curMode = normal
		if e.b.cursor.ofset > 0 {
			e.b.cursor.ofset -= 1
		}
	case '\127', '\x7f':
		e.b.RemoveKey(0)
	case '\t':
		//NOTE: yeah, I just insert 4 spaces instead of tabs
		e.b.InsertKey(' ')
		e.b.InsertKey(' ')
		e.b.InsertKey(' ')
		e.b.InsertKey(' ')
	default:
		fmt.Print(key)
		e.b.InsertKey(key)
	}
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
		fmt.Print(clearView, clearHistory, moveToStart)
		term.Restore(e.fdIn, e.oldState)
		os.Exit(0)
		return true
	case "w":
		err := e.SaveFile()
		if err != nil {
			return false
		}
		return true
	case "x", "wq":
		err := e.SaveFile()
		if err != nil {
			return false
		} else {
			fmt.Print(clearView, clearHistory, moveToStart)
			term.Restore(e.fdIn, e.oldState)
			os.Exit(0)
			return true
		}
	default:
		return false
	}
}

func (e *Editor) Run() {
	//var buf [1]rune
	Draw(e)
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
		}

		Draw(e)
	}
}
