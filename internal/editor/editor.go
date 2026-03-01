package editor

import (
	"Enot/Bard/config"
	tui "Enot/Bard/internal/TUI"
	"Enot/Bard/internal/ascii"
	"Enot/Bard/internal/buffer"
	"Enot/Bard/internal/mode"
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Editor struct {
	oldState *term.State
	b        *buffer.Buffer
	tui      *tui.TUI
	c        *config.Config
	curMode  mode.Mode
	command  string // used in command mode
	subCmd   string // sub command, line 12j
	file     string
	save     bool // is terminal save to work (depends on window size)
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

	_b := buffer.InitBuffer()
	_tui := tui.InitTUI(_h, _w)
	_c := config.InitConfig()

	e := &Editor{
		oldState: old,
		b:        _b,
		tui:      _tui,
		c:        _c,
		curMode:  mode.Normal,
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

	return e
}

func (e *Editor) TermSizeMonitor() {
	e.tui.TermSizeMonitor(e.fdOut)
}

func (e *Editor) Exit(code int) {
	e.c.Save()

	fmt.Print(ascii.ClearView, ascii.ClearHistory, ascii.MoveToStart, ascii.CursorReset, ascii.ResetTerminal)
	term.Restore(e.fdIn, e.oldState)
	os.Exit(code)
}

// Main loop
func (e *Editor) Run() {
	fmt.Print(ascii.SaveTerminal)
	e.Draw()
	reader := bufio.NewReader(os.Stdin)
	for {
		e.tui.Message = ""
		key, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}
		switch e.curMode {
		case mode.Normal:
			e.caseNormal(key)
		case mode.Command:
			e.caseCommand(key)
		case mode.Insert:
			e.caseInsert(key)
		case mode.Visual:
			e.caseVisual(key)
		case mode.Visual_line:
			e.caseVisualLine(key)
		case mode.Replace:
			e.caseReplaceMode(key)
		default:
			e.Exit(1)
		}

		e.setUiCursor()
		e.Draw()
	}
}
