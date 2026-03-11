package main

import (
	"os"

	"github.com/EnotInc/Bard/internal/editor"
)

func main() {
	var file string
	e := editor.InitEditor()
	if len(os.Args) == 2 {
		file = os.Args[1]
		e.LoadFile(file)
	}

	e.TermSizeMonitor()
	e.Run()
}
