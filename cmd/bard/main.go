package main

import (
	"os"

	"github.com/EnotInc/Bard/internal/editor"
)

func main() {
	e := editor.InitEditor()
	if len(os.Args) == 2 {
		arg := os.Args[1]

		switch arg {
		case "-h", "--help":
			e.StartHelp()
		default:
			e.LoadFile(arg)
		}
	}

	e.TermSizeMonitor()
	e.Run()
}
