package editor

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

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

func (e *Editor) execCommand() {
	switch e.command {
	case "q":
		fmt.Print(clearView, clearHistory, moveToStart, cursorReset)
		term.Restore(e.fdIn, e.oldState)
		os.Exit(0)
	case "w":
		e.SaveFile()
	case "x", "wq":
		e.SaveFile()
		fmt.Print(clearView, clearHistory, moveToStart, cursorReset)
		term.Restore(e.fdIn, e.oldState)
		os.Exit(0)
	case "rln":
		e.ui.rln = !e.ui.rln
	default:
		e.message = "unknown command"
	}
}
