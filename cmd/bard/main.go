package main

import (
	"Enot/Bard/editor"
	"os"
)

func main() {
	var file string
	e := editor.InitEditor()
	if len(os.Args) == 2 {
		file = os.Args[1]
		e.LoadFile(file)
	}

	go e.TermSizeMonitor()
	e.Run()
}
