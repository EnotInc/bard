package editor

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func (e *Editor) caseCommand(key rune) {
	switch key {
	case '\033':
		e.curCommand = ""
		e.curMode = normal
	case '\127', '\x7f':
		if len(e.curCommand) > 0 {
			e.curCommand = e.curCommand[:len(e.curCommand)-1]
		} else {
			e.curCommand = ""
			e.curMode = normal
		}
	case '\013', '\r', '\n':
		e.execCommand()
		e.curCommand = ""
		e.curMode = normal
	default:
		e.curCommand += string(key)
	}
}

func (e *Editor) execCommand() {
	switch e.curCommand {
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
