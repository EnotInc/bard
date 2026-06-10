package main

import (
	"os"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/editor"
	"github.com/EnotInc/Bard/internal/explorer"
	"github.com/EnotInc/Bard/internal/screen"
)

func main() {
	config.InitConfig()
	screen.InitScreen()

	border := true
	h := screen.H()
	ex_w := 30
	ed_w := screen.W() - ex_w

	ed := editor.InitEditor(ed_w, h)
	ed_tile, err := screen.NewTile(ed, ed_w, h, border)
	if err != nil {
		panic(err)
	}

	ex := explorer.InitExplorer(ed.OpenFileCallback, ed.RemoveFileCallback, ex_w, h)
	ex_tile, err := screen.NewTile(ex, ex_w, h, border)
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 {
		arg := os.Args[1]

		switch arg {
		case "-h", "--help":
			ed.StartHelp()
		default:
			ed.LoadFile(arg)
		}
	}

	screen.SetStatusBar(ed.DrawStatusBar)
	screen.AddTile(ed_tile)
	screen.AddTile(ex_tile)
	screen.HideTile()

	screen.TermSizeMonitor()

	screen.Run()
}
