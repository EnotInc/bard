package main

import (
	"os"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/popups"
	"github.com/EnotInc/Bard/internal/popups/setting"
	"github.com/EnotInc/Bard/internal/popups/themes"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/tiles/editor"
	"github.com/EnotInc/Bard/internal/tiles/explorer"
	"github.com/EnotInc/Bard/theme"
)

func main() {
	config.InitConfig()
	theme.InitTheme(config.GetConfig().ThemeName)
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
		ed.SetErrorCallback,
		ex_w, h)

	ex_tile, err := screen.NewTile(ex, 0.25)
	if err != nil {
		panic(err)
	}

	screen.SetStatusBar(ed.DrawStatusBar)
	screen.AddTile(ed_tile)
	screen.AddTile(ex_tile)

	t := themes.InitThemes(
		ed.PurgeCacheCallback,
		ed.ChangeModeCallback,
		ed.SetErrorCallback)
	t_popup := screen.NewPopup(t)

	s := setting.IntiSettings(ed.PurgeCacheCallback)
	s_popup := screen.NewPopup(s)

	screen.AddPopup(t_popup, popups.Themes)
	screen.AddPopup(s_popup, popups.Settings)

	if len(os.Args) == 2 {
		arg := os.Args[1]

		switch arg {
		case "-h", "--help":
			ed.StartHelp()
			screen.HideTile()
		case "--space", "-s":
			space := config.GetSpacePath()
			screen.SetRoot(space)
			ex.SetPath(space)
		default:
			_, err := os.Stat(arg)
			if err != nil {
				ed.CreateFile(arg)
			}
			f, err := os.Stat(arg)
			if err != nil {
				panic(err)
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
