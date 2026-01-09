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
	case 'j':
	case 'k':
	case 'l':
	case 'i':
		e.curMode = insert
	case 'a':
	case ':':
		e.curMode = command
	case 'o':
	case 'O':
	}
}

func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b.InsertLine()
	case '\033':
		e.curMode = normal
	case '\127', '\x7f':
	default:
		fmt.Print(key)
		e.b.InsertKey(key)
	}
}

func (e *Editor) caseCommand(key rune) {
	switch key {
	case '\033':
	case '\127', '\x7f':
	case '\013', '\r', '\n':
		success := e.execCommand()
		if !success {
			e.curCommand = ""
			e.curMode = normal
			//TODO: move cursor back to prev pos
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
