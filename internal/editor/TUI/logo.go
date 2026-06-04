package tui

import "github.com/EnotInc/Bard/internal/enums/ascii"

const (
	hello   = "Oh, hello there!"
	info    = ":q - quit        :w <file name?> - save        :x - save and quit"
	motions = "If you didn't know this, well, it's time to learn some vim motions"
)

// Called when used starts emtpy editor
// used to get ascii art ot display it in the middle of screen
// result will be saved at tui.Hello field
func (tui *TUI) BuidASCII() {

	data := [][]rune{}
	data = append(data, []rune(""))
	if tui.Save {
		data = append(data, ascii.Bard...)
	}
	data = append(data, []rune(""))
	data = append(data, []rune(hello))
	data = append(data, []rune(""))
	data = append(data, []rune(info))
	data = append(data, []rune(motions))

	tui.Hello = data
}

// Called to get one line at the time from tui.Hello (by given index)
func (tui *TUI) GetASCIIInfo(index int) []rune {

	if index >= len(tui.Hello) {
		return []rune{}
	}

	return tui.Hello[index]
}
