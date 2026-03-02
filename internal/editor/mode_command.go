package editor

import "Enot/Bard/internal/mode"

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
		if len(e.b) > 1 {
			e.delBuffer(e.curBuffer)
		} else {
			e.SaveFile()
			e.Exit(0)
		}
	case "rln":
		e.c.RLN = !e.c.RLN
	case "showmd":
		e.c.ShowMD = !e.c.ShowMD
	case "render", "rnd":
		e.c.Render = !e.c.Render
	case "help", "h":
		e.newBuffer()
		e.LoadFile("docs/help/help.md")
	default:
		if len(e.command) > 3 {
			if e.command[0] == 'w' && e.command[1] == ' ' {
				fileName := e.command[2:]
				e.CreateFile(fileName)
				e.b[e.curBuffer].Title = fileName
				e.SaveFile()
			} else {
				e.tui.Message = "unknown command"
			}
			return
		}
		e.tui.Message = "unknown command"
	}
}
