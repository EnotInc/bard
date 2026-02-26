package editor

func (e *Editor) caseCommand(key rune) {
	switch key {
	case '\033':
		e.command = ""
		e.curMode = normal
	case '\127', '\x7f':
		if len(e.command) > 0 {
			e.command = e.command[:len(e.command)-1]
		} else {
			e.command = ""
			e.curMode = normal
		}
	case '\013', '\r', '\n':
		e.execCommand()
		e.command = ""
		e.curMode = normal
	default:
		e.command += string(key)
	}
}

// For now I just compare commands, and run them
// Later I'll make some sort of a lexer to do it
func (e *Editor) execCommand() {
	switch e.command {
	case "q":
		e.Exit(0)
	case "w":
		e.SaveFile()
	case "x", "wq":
		e.SaveFile()
		e.Exit(0)
	case "rln":
		e.c.RLN = !e.c.RLN
	case "showmd":
		e.c.ShowMD = !e.c.ShowMD
	case "render", "rnd":
		e.c.Render = !e.c.Render
	default:
		if len(e.command) > 3 {
			if e.command[0] == 'w' && e.command[1] == ' ' {
				fileName := e.command[2:]
				e.CreateFile(fileName)
				e.file = fileName
				e.SaveFile()
			} else {
				e.message = "unknown command"
			}
			return
		}
		e.message = "unknown command"
	}
}
