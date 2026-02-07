package editor

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
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
	message    string
	isMdFile   bool
	fdIn       int
}

func InitEditor() *Editor {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil {
		panic(err)
	}
	_w, _h, _ := term.GetSize(_fdOut)
	if _w < 80 || _h < 20 { //standart terminal size
		panic(fmt.Sprintf("\nUnable to run in this terminal, window is too small: (%d %d)", _w, _h))
	}

	_b := InitBuffer()
	_ui := InitUI(_h, _w)

	e := &Editor{
		oldState:   old,
		b:          _b,
		ui:         _ui,
		curMode:    normal,
		isMdFile:   false,
		curCommand: "",
		fdIn:       _fdIn,
	}

	return e
}

func (e *Editor) Run() {
	e.ui.Draw(e)
	reader := bufio.NewReader(os.Stdin)
	for {
		e.message = ""
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
