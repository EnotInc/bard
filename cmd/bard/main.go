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
		if _, err := os.Stat(file); err == nil {
			e.LoadFile(file)
		} else {
			os.Create(file)
		}
	}
	e.Run()
}
