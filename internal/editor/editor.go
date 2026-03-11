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
	"strconv"

	"golang.org/x/term"
)

type Editor struct {
	oldState  *term.State
	b         []*buffer.Buffer //list of buffers
	tui       *tui.TUI
	c         *config.Config
	curMode   mode.Mode
	command   string // used in command mode
	subCmd    string // sub command, line 12j
	save      bool   // is terminal save to work (depends on window size)
	fdOut     int
	fdIn      int
	curBuffer int // current buffer index
}

func InitEditor() *Editor {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil {
		panic(err)
	}
	_w, _h, _ := term.GetSize(_fdOut)

	_c := config.InitConfig()
	_b := buffer.InitBuffer()
	_tui := tui.InitTUI(_h, _w)

	e := &Editor{
		oldState:  old,
		b:         _b,
		tui:       _tui,
		c:         _c,
		curMode:   mode.Normal,
		command:   "",
		subCmd:    "",
		fdOut:     _fdOut,
		fdIn:      _fdIn,
		curBuffer: 0,
	}

	if _w < 80 || _h < 30 {
		e.tui.Save = false
	}

	e.tui.BuidASCII()
	return e
}

func (e *Editor) TermSizeMonitor() {
	go e.tui.TermSizeMonitor(e.fdOut)
	go e.listenResize()
}

func (e *Editor) listenResize() {
	for {
		value := <-e.tui.Redraw
		if value {
			e.Draw()
		}
	}
}

func (e *Editor) Exit(code int) {
	e.c.Save()

	fmt.Print(ascii.ClearView, ascii.ClearHistory, ascii.MoveToStart, ascii.CursorReset, ascii.ResetTerminal)
	term.Restore(e.fdIn, e.oldState)
	os.Exit(code)
}

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
