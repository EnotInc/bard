package editor

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	clearView    = "\033[2J"
	clearHistory = "\033[3J"
	moveToStart  = "\033[0H"
	cursorReset  = "\033]112\a"

	saveTerminal  = "\033[?1049h"
	resetTerminal = "\033[?1049l"
)

type Mode string

const (
	normal      Mode = "NORMAL"
	command     Mode = "COMMAND"
	insert      Mode = "INSERT"
	visual      Mode = "VISUAL"
	visual_line Mode = "VISUAL LINE"
)

type Editor struct {
	oldState  *term.State
	b         *Buffer
	ui        *UI
	c         *Config
	curMode   Mode
	command   string
	subCmd    string
	file      string
	message   string
	showHello bool
	save      bool
	isMdFile  bool
	fdOut     int
	fdIn      int
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
	_c := InitConfig()

	e := &Editor{
		oldState: old,
		b:        _b,
		ui:       _ui,
		c:        _c,
		curMode:  normal,
		isMdFile: false,
		save:     true,
		command:  "",
		subCmd:   "",
		fdOut:    _fdOut,
		fdIn:     _fdIn,
	}

	if _w < 80 || _h < 30 { // standard terminal size
		e.save = false
	}
	e.buidASCII()

	return e
}

// This function is called in the main.go file in a goroutine.
// Here I just recalculate the terminal size and adjust Bard to it
func (e *Editor) TermSizeMonitor() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var last_w, last_h = e.ui.w, e.ui.h

	for range ticker.C {
		w, h, err := term.GetSize(e.fdOut)
		if err != nil {
			continue
		}

		if last_w != w || last_h != h {
			last_w = w
			last_h = h

			e.resize(w, h)
			e.ui.Draw(e)
		}
	}
}

func (e *Editor) resize(w int, h int) {
	e.ui.w = w
	e.ui.h = h
	e.setUiCursor()
}

func (e *Editor) Exit(code int) {
	e.c.Save()

	fmt.Print(clearView, clearHistory, moveToStart, cursorReset, resetTerminal)
	term.Restore(e.fdIn, e.oldState)
	os.Exit(code)
}

// Main loop
func (e *Editor) Run() {
	fmt.Print(saveTerminal)
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
		case visual_line:
			e.caseVisualLine(key)
		}

		e.setUiCursor()
		e.ui.Draw(e)
	}
}
