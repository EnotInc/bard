package main

import (
	"os"

	"github.com/EnotInc/Bard/internal/editor"
	"github.com/EnotInc/Bard/internal/screen"
)

func main() {
	// e.TermSizeMonitor()
	screen.InitScreen()

	border := true
	w := screen.W()
	h := screen.H()

	e := editor.InitEditor(w, h)
	e_tile, err := screen.NewTile(e, w, h, 0, 0, border)
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 {
		arg := os.Args[1]

		switch arg {
		case "-h", "--help":
			e.StartHelp()
		default:
			e.LoadFile(arg)
		}
	}

	screen.SetStatusBar(e.DrawStatusBar)
	screen.AddTile(e_tile)

	screen.Run()
}
