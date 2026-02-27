package tui

import "Enot/Bard/internal/ascii"

const (
	hello   = "Oh, hello there!"
	info    = ":q - quit        :w <file name?> - save        :x - save and quit"
	motions = "If you didn't know this, well, it's time to learn some vim motions"
)

func (tui *TUI) buidASCII() {

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

func (tui *TUI) GetASCIIInfo(index int) []rune {

	if index >= len(tui.Hello) {
		return []rune{}
	}

	return tui.Hello[index]
}
