package editor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/buffer"

	tui "github.com/EnotInc/Bard/internal/TUI"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

// This is main structure, that contains current editor state
// oldState - used to work with raw terminal mode
// b - list of Buffer. List is used in work with tabs
// tui - TUI
// CurMode - current editor Mode
// hash - used in DrawDiff() func for diff rendering
// command - used in command mode, stores user input
// subCmd - sub command. Used to store commands like `12k`
// save - is terminal save to work in (depends on window size, if w < 80 or h < 30 then terminal is not save)
// fdOut - used to work with raw input
// fdIn - used to work with raw input
// curBuffer - current buffer index
type Editor struct {
	oldState  *term.State
	b         []*buffer.Buffer
	tui       *tui.TUI
	theme     *config.Theme
	curMode   mode.Mode
	hash      map[int]uint32
	command   string
	subCmd    string
	lastCmd   string
	IsChanged bool
	save      bool
	fdOut     int
	fdIn      int
	curBuffer int
}

// Initialisation of editor
// turn terminal into raw mode, saves old state, initializes Config, Buffer and TUI
// checks if terminal save to work in
func InitEditor() *Editor {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil {
		panic(err)
	}
	_w, _h, _ := term.GetSize(_fdOut)
	if _w <= 40 || _h < 15 {
		panic("Unable to run Bard. Window size is too small!")
	}

	config.InitConfig()
	cfg := config.Get()
	_t, err := config.InitTheme(cfg.ThemeName)
	_b := buffer.InitBuffer()
	_tui := tui.InitTUI(_h, _w, _t)

	e := &Editor{
		oldState:  old,
		b:         _b,
		tui:       _tui,
		theme:     _t,
		curMode:   mode.Normal,
		hash:      make(map[int]uint32),
		command:   "",
		subCmd:    "",
		fdOut:     _fdOut,
		fdIn:      _fdIn,
		curBuffer: 0,
	}

	if err != nil {
		cfg.ThemeName = cfg.DefaultThemeName()
	}

	if _w < 80 || _h < 30 {
		e.tui.Save = false
	}

	e.tui.BuidASCII()
	return e
}

// Starts 2 goroutines to check if terminal window size was changed:
// e.tui.TermSizeMotitor and e.listenResize
func (e *Editor) TermSizeMonitor() {
	go e.tui.TermSizeMonitor(e.fdOut)
	go e.listenResize()
}

// wait till e.tui.Redraw is true, and then it Draw() editor with new size
func (e *Editor) listenResize() {
	for {
		value := <-e.tui.Redraw
		if value {
			e.DrawDiff()
		}
	}
}

// TODO: add message
// Used to restore old terminal state, change terminal buffer (via ascii escape sequence) and stop Bard with status code
func (e *Editor) Exit(code int) {
	config.Save()

	fmt.Print(ascii.ClearView, ascii.ClearHistory, ascii.MoveToStart, ascii.CursorReset, ascii.ResetTerminal, ascii.ResetCursor)
	term.Restore(e.fdIn, e.oldState)
	if r := recover(); r != nil {
		err := e.saveLog(r)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Bard stopped with error. More information you can find in '~/.bard/.log' file")
		}
	}
	os.Exit(code)
}

// Used to get move func (one of 'H', 'J', 'K' or 'L'), and move cursor by some amout
// needded about of moves is get's from subCmd, and if it not parsed, it does nothing
func (e *Editor) execWithSubCommand(exec func(int)) {
	if e.subCmd == "" {
		exec(1)
		return
	}
	amount, err := strconv.Atoi(e.subCmd)
	if err != nil {
		e.subCmd = ""
		return
	}
	exec(amount)
	e.subCmd = ""
}

// FIXME: does not work...
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

// Gets user input, switched by current mode to decide what to do with pressed key and calles DrawDiff() to display changes
func (e *Editor) Run() {
	defer e.Exit(1)
	fmt.Print(ascii.SaveTerminal, ascii.ClearView, ascii.ClearHistory)
	e.DrawDiff()
	reader := bufio.NewReader(os.Stdin)
	for {
		e.tui.Message = ""
		key, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}

		switch e.curMode {
		case mode.Normal:
			if e.IsGeneralMove(key) {
				e.GeneralCase(key)
			} else {
				e.caseNormal(key)
			}
		case mode.Visual:
			if e.IsGeneralMove(key) {
				e.GeneralCase(key)
			} else {
				e.caseVisual(key)
			}
		case mode.Visual_line:
			if e.IsGeneralMove(key) {
				e.GeneralCase(key)
			} else {
				e.caseVisualLine(key)
			}
		case mode.Command:
			e.caseCommand(key)
		case mode.Insert:
			e.caseInsert(key)
		case mode.Replace:
			e.caseReplaceMode(key)
		default:
			e.Exit(1)
		}

		e.setUiCursor()
		e.DrawDiff()
	}
}
