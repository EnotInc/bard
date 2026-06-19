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
	config.CreateSpace()
	screen.InitScreen()

	h := screen.H()
	ex_w := 30
	ed_w := screen.W() - ex_w

	ed := editor.InitEditor(ed_w)
	ed_tile, err := screen.NewTile(ed, 0.75)
	if err != nil {
		panic(err)
	}

	ex := explorer.InitExplorer(
		ed.OpenFileCallback,
		ed.RemoveFileCallback,
		ed.RenameCallback,
		ed.ChangeModeCallback,
		ex_w, h)

	ex_tile, err := screen.NewTile(ex, 0.25)
	if err != nil {
		panic(err)
	}

	screen.SetStatusBar(ed.DrawStatusBar)
	screen.AddTile(ed_tile)
	screen.AddTile(ex_tile)

	if len(os.Args) == 2 {
		arg := os.Args[1]

		switch arg {
		case "-h", "--help":
			ed.StartHelp()
		case "--space", "-s":
			space := config.GetSpacePath()
			screen.SetRoot(space)
			ex.SetPath(space)
		default:
			f, err := os.Stat(arg)
			if err != nil {
				ed.CreateFile(arg)
			}

			if !f.IsDir() {
				ed.LoadFile(arg)
				screen.HideTile()
			} else {
				screen.SetRoot(arg)
				ex.SetPath(arg)
			}
		}
	}
	screen.ShiftFocus()

	screen.TermSizeMonitor()

	screen.Run()
}
