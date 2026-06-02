package editor

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/help"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

// Called from Run() func
// Used to decide what do to with pressed key
func (e *Editor) caseCommand(key rune) {
	switch key {
	case keys.Esc:
		e.cmd.command = ""
		e.curMode = mode.Normal

	case keys.Backspace:
		if len(e.cmd.command) > 0 {
			e.cmd.command = e.cmd.command[:len(e.cmd.command)-1]
		} else {
			e.cmd.command = ""
			e.curMode = mode.Normal
		}

	case keys.Enter:
		e.execCommand()
		e.cmd.command = ""
		e.curMode = mode.Normal
	case '>':
		if e.cmd.index < len(e.cmd.history) {
			e.cmd.command = e.cmd.history[e.cmd.index]
			e.cmd.index += 1
		}
	case '<':
		if e.cmd.index > 0 {
			e.cmd.index -= 1
			e.cmd.command = e.cmd.history[e.cmd.index]
		}

	default:
		if unicode.IsPrint(key) {
			e.cmd.command += string(key)
		}
	}
}

func (cmd *cmd) saveToHisory() {
	// if we moved back on cmd history and change one of the commands, history will continue from this point
	if cmd.index != len(cmd.history) {
		cmd.history = slices.Delete(cmd.history, cmd.index, len(cmd.history))
	}
	cmd.history = append(cmd.history, cmd.command)
	cmd.index = len(cmd.history)
}

// For now I just compare commands, and run them
// Later I'll make some sort of a lexer to do it
func (e *Editor) execCommand() {
	cfg := config.Get()
	e.cmd.saveToHisory()
	switch e.cmd.command {
	case "q":
		if len(e.b) > 1 {
			e.delBuffer(e.curBuffer)
		} else {
			e.Exit(0)
		}

	case "qa":
		e.Exit(0)

	case "w":
		e.SaveFile()

	case "x", "wq":
		e.SaveFile()
		if len(e.b) > 1 {
			e.delBuffer(e.curBuffer)
		} else {
			e.Exit(0)
		}

	case "help", "h":
		e.OpenHelp(help.About)

	case "rln":
		cfg.RLN = !cfg.RLN

	case "showmd":
		cfg.ShowMD = !cfg.ShowMD
		e.IsChanged = true
		e.tui.PurgeCache()
		e.PurgeCache()

	case "render", "rnd":
		cfg.Render = !cfg.Render
		e.IsChanged = true
		e.tui.PurgeCache()
		e.PurgeCache()

	case "tn", "tabnames":
		cfg.TabNames = !cfg.TabNames

	case "gt":
		e.nextTab()

	case "gT":
		e.prevTab()

	case "newtab", "nt":
		e.newBuffer()

	case "theme":
		e.tui.Message = fmt.Sprintf("Theme: %s", cfg.ThemeName)

	default:
		e.parseCommand()
	}
}

// Used to parse some specific commands like `help`, or `w` (save)
func (e *Editor) parseCommand() {
	if len(e.cmd.command) >= 3 {
		parts := strings.Split(e.cmd.command, " ")
		if len(parts) != 2 {
			e.tui.Message = "bad syntax"
			return
		}

		cfg := config.Get()

		cmd := parts[0]
		arg := parts[1]

		switch cmd {
		case "w":
			e.CreateFile(arg)
			e.b[e.curBuffer].Title = arg
			e.SaveFile()

		case "newtab", "nt":
			e.newBuffer()
			e.LoadFile(arg)
			e.b[e.curBuffer].Title = arg

		case "help", "h":
			var topic help.Topic = help.Topic(arg)
			e.OpenHelp(topic)

		case "theme":
			if arg == "reload" {
				arg = cfg.ThemeName
			}

			msg := e.theme.ChangeTheme(arg)
			if msg != "" {
				e.tui.Message = msg
				return
			}
			cfg.ThemeName = arg
			e.tui.PurgeCache()
			e.PurgeCache()

		case "gt":
			page, err := strconv.Atoi(arg)
			if err != nil {
				e.tui.Message = "unable to get tab number"
				return
			}

			page -= 1
			if page < 0 || page > len(e.b) {
				e.tui.Message = "can't open this tab"
				return
			}

			e.curBuffer = page
		case "tabstop", "ts":
			ts, err := strconv.Atoi(arg)
			if err != nil {
				e.tui.Message = "unable to get tabstop"
				return
			}

			cfg.TabStop = ts
			config.FixConfig()

			e.tui.PurgeCache()
			e.PurgeCache()

		default:
			e.tui.Message = "unknown command"
		}
	}
}
