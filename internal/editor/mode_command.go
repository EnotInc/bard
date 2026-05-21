package editor

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/EnotInc/Bard/internal/enums"
)

// Called from Run() func
// Used to decide what do to with pressed key
func (e *Editor) caseCommand(key rune) {
	switch key {
	case '\033':
		e.command = ""
		e.curMode = enums.Normal
	case '\127', '\x7f':
		if len(e.command) > 0 {
			e.command = e.command[:len(e.command)-1]
		} else {
			e.command = ""
			e.curMode = enums.Normal
		}
	case '\013', '\r', '\n':
		e.execCommand()
		e.command = ""
		e.curMode = enums.Normal
	default:
		if unicode.IsPrint(key) {
			e.command += string(key)
		}
	}
}

// For now I just compare commands, and run them
// Later I'll make some sort of a lexer to do it
func (e *Editor) execCommand() {
	switch e.command {
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
		e.OpenHelp(enums.HelpAbout)
	case "rln":
		e.c.RLN = !e.c.RLN
	case "showmd":
		e.c.ShowMD = !e.c.ShowMD
		e.IsChanged = true
	case "render", "rnd":
		e.c.Render = !e.c.Render
		e.IsChanged = true
	case "tn", "tabnames":
		e.c.TabNames = !e.c.TabNames
	case "gt":
		e.nextTab()
	case "gT":
		e.prevTab()
	case "newtab", "nt":
		e.newBuffer()
	case "theme":
		e.tui.Message = fmt.Sprintf("Theme: %s", e.c.ThemeName)
	default:
		e.parseCommand()
	}
}

// Used to parse some specific commands like `help`, or `w` (save)
func (e *Editor) parseCommand() {
	if len(e.command) >= 3 {
		parts := strings.Split(e.command, " ")
		if len(parts) != 2 {
			e.tui.Message = "bad syntax"
			return
		}

		cmd := parts[0]
		arg := parts[1]

		switch cmd {
		case "w":
			e.CreateFile(arg)
			e.b[e.curBuffer].Title = arg
			e.SaveFile()
		case "newtab", "nt":
			e.newBuffer()
			e.CreateFile(arg)
			e.b[e.curBuffer].Title = arg
		case "help", "h":
			var topic enums.Help = enums.Help(arg)
			e.OpenHelp(topic)
		case "theme":
			msg := e.theme.ChangeTheme(arg)
			if msg != "" {
				e.tui.Message = msg
				return
			}
			e.tui.PurgeCache()
			e.c.ThemeName = arg
			e.c.Save()
		default:
			e.tui.Message = "unknown command"
		}
	}
}
