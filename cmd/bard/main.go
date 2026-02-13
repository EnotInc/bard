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
		if _, err := os.Stat(file); err != nil {
			os.Create(file)
		}
		e.LoadFile(file)
	}

	go e.TermSizeMonitor()
	e.Run()
}
