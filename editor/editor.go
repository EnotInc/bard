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
	oldState *term.State
	b        *Buffer
	ui       *UI
	curMode  Mode
	command  string
	subCmd   string
	file     string
	message  string
	isMdFile bool
	fdOut    int
	fdIn     int
}

func InitEditor() *Editor {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil {
		panic(err)
	}
	_w, _h, _ := term.GetSize(_fdOut)
	if _w < 80 || _h < 20 { // standard terminal size
		panic(fmt.Sprintf("\nUnable to run in this terminal, window is too small: (%d %d)", _w, _h))
	}

	_b := InitBuffer()
	_ui := InitUI(_h, _w)

	e := &Editor{
		oldState: old,
		b:        _b,
		ui:       _ui,
		curMode:  normal,
		isMdFile: false,
		command:  "",
		subCmd:   "",
		fdOut:    _fdOut,
		fdIn:     _fdIn,
	}

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

// Main loop
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
		case visual_line:
			e.caseVisualLine(key)
		}

		e.setUiCursor()
		e.ui.Draw(e)
	}
}
