package editor

import (
	"strconv"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/editor/buffer"

	tui "github.com/EnotInc/Bard/internal/editor/TUI"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

// This is main structure, that contains current editor state
// b - list of Buffer. List is used in work with tabs
// tui - TUI
// CurMode - current editor Mode
// command - used in command mode, stores user input
// subCmd - sub command. Used to store commands like `12k`
// save - is terminal save to work in (depends on window size, if w < 80 or h < 30 then terminal is not save)
// curBuffer - current buffer index
type Editor struct {
	tui             *tui.TUI
	cmd             *cmd
	curMode         mode.Mode
	emptyLineSpaces string
	subCmd          string
	lastCmd         string
	b               []*buffer.Buffer
	curBuffer       int
	save            bool
}

type cmd struct {
	command []rune
	history [][]rune
	index   int
}

func initCmd() *cmd {
	return &cmd{
		command: []rune{},
		history: [][]rune{},
		index:   0,
	}
}

// Initialisation of editor
// turn terminal into raw mode, saves old state, initializes Config, Buffer and TUI
// checks if terminal save to work in
func InitEditor(w int) *Editor {
	cfg := config.GetConfig()
	err := config.InitTheme(cfg.ThemeName)
	_b := buffer.InitBuffer()
	_tui := tui.InitTUI(w)
	_cmd := initCmd()

	e := &Editor{
		b:         _b,
		tui:       _tui,
		cmd:       _cmd,
		curMode:   mode.Normal,
		lastCmd:   "",
		subCmd:    "",
		curBuffer: 0,
	}

	if err != nil {
		cfg.ThemeName = cfg.DefaultThemeName()
	}

	if w < 80 {
		e.tui.Save = false
	}

	return e
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

func (e *Editor) SetErrorCallback(err string) {
	e.tui.Error = err
}
