package main

import (
	"Enot/Bard/editor"
	"os"
)

func main() {
	e := editor.InitEditor()
	if len(os.Args) == 2 {
		filePath := os.Args[1]
		err := e.LoadFile(filePath)
		if err != nil {
			panic(err)
		}
	}
	e.Run()
}
