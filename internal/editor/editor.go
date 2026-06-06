package editor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/editor/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/screen"

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
	emtpyLineSpases string
	subCmd          string
	lastCmd         string
	b               []*buffer.Buffer
	curBuffer       int
	IsChanged       bool
	save            bool
}

type cmd struct {
	command string
	history []string
	index   int
}

func initCmd() *cmd {
	return &cmd{
		command: "",
		history: []string{},
		index:   0,
	}
}

// Initialisation of editor
// turn terminal into raw mode, saves old state, initializes Config, Buffer and TUI
// checks if terminal save to work in
func InitEditor(w, h int) *Editor {
	config.InitConfig()
	cfg := config.GetConfig()
	err := config.InitTheme(cfg.ThemeName)
	_b := buffer.InitBuffer()
	_tui := tui.InitTUI(h, w)
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

	if w < 80 || h < 30 {
		e.tui.Save = false
	}

	e.tui.BuidASCII()
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

func (e *Editor) Handle(key rune) {
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
		screen.Exit(1)
	}

	e.setUiCursor()
}

func (e *Editor) GetCursor(withBorder bool) (int, int) {
	var x int
	var y int
	if e.curMode == mode.Command {
		x = len(e.cmd.command) + len(e.emtpyLineSpases)
		y = e.tui.H

		if !withBorder {
			x += 1
		}
	} else {
		x = e.tui.CurOff + enums.InitialOffset + len(e.emtpyLineSpases)
		y = e.tui.CurRow + enums.CursorOffset
	}
	return x, y
}

func (e *Editor) SetTitle() string {
	var tabs []string
	for _, t := range e.b {
		tabs = append(tabs, t.Title)
	}
	cfg := config.GetConfig()
	return e.tui.BuildTabs(tabs, e.curBuffer, cfg.TabNames)
}

func (e *Editor) PreDraw() {
	e.emtpyLineSpases = tui.BuildSpaces(len(fmt.Sprint(len(e.b[e.curBuffer].Lines))))
	for i := range e.tui.YScroll {
		curLine := string(e.b[e.curBuffer].Lines[i].Data)
		if strings.HasPrefix(curLine, "```") {
			e.tui.ToggleRender()
		}
	}
}
