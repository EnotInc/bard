package editor

import (
	"strconv"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/editor/buffer"

	tui "github.com/EnotInc/Bard/internal/editor/TUI"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

type Editor struct {
	tui             *tui.TUI
	cmd             *cmd
	curMode         mode.Mode
	emptyLineSpaces string
	subCmd          string
	lastCmd         string
	b               []*buffer.Buffer
	curBuffer       int
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

func (e *Editor) replaceWithAmount(key rune) {
	if e.subCmd == "r" {
		e.caseReplaceChar(key)
		return
	}
}

func (e *Editor) SetErrorCallback(err string) {
	e.tui.Error = err
}
