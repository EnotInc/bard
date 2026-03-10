package editor

import (
	"Enot/Bard/internal/enums"
	"Enot/Bard/internal/mode"
	"strings"
)

func (e *Editor) caseCommand(key rune) {
	switch key {
	case '\033':
		e.command = ""
		e.curMode = mode.Normal
	case '\127', '\x7f':
		if len(e.command) > 0 {
			e.command = e.command[:len(e.command)-1]
		} else {
			e.command = ""
			e.curMode = mode.Normal
		}
	case '\013', '\r', '\n':
		e.execCommand()
		e.command = ""
		e.curMode = mode.Normal
	default:
		e.command += string(key)
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
		e.OpenHelp(enums.About)
	case "rln":
		e.c.RLN = !e.c.RLN
	case "showmd":
		e.c.ShowMD = !e.c.ShowMD
	case "render", "rnd":
		e.c.Render = !e.c.Render
	case "tn", "tabnames":
		e.c.TabNames = !e.c.TabNames
	case "gt":
		e.nextTab()
	case "gT":
		e.prevTab()
	case "newtab", "nt":
		e.newBuffer()
	default:
		e.parceCommand()
	}
}

func (e *Editor) parceCommand() {
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
			e.tui.Message = "not implemented yet"
			var topic enums.Help = enums.Help(arg)
			e.OpenHelp(topic)
		default:
			e.tui.Message = "unknown command"
		}
	}
}
