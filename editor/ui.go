package editor

import (
	"fmt"
	"strconv"
)

type color string

const (
	reset  color = "\033[0m"
	red    color = ""
	green  color = ""
	gray   color = ""
	yellow color = "\033[33m"
	black  color = ""
)

const (
	initialOfset    = 3
	cursorLineOfset = 1
)

type UI struct {
	upperBorder int
	lowerBorder int
	curRow      int
}

func InitUI(h int) *UI {
	ui := &UI{
		upperBorder: 0,
		lowerBorder: h,
		curRow:      0,
	}
	return ui
}

func colorise(data string, c color) string {
	return fmt.Sprintf("%s%s%s", c, data, reset)
}

func (b *Buffer) buildNumber(n int, maxOfset int) string {
	numStr := strconv.Itoa(n)
	numLen := len(numStr)
	num := ""
	if maxOfset <= 4 {
		maxOfset = 4
	}
	for range maxOfset - numLen {
		num += " "
	}
	num += numStr
	if b.cursor.line+1 == n {
		num = colorise(num, yellow)
	} else {
		num = colorise(num, reset)
	}
	return fmt.Sprintf("%s", num)
}

func buildSpaces(maxOfset int) string {
	space := ""
	if maxOfset <= 4 {
		maxOfset = 4
	}
	for range maxOfset - 1 {
		space += " "
	}
	return fmt.Sprintf("%s", space)
}

// func (ui *UI) Draw(e *Editor) {
// 	emtpyLineSpases := buildSpaces(len(fmt.Sprint(len(e.b.lines))))
// 	data := ""
// 	data += fmt.Sprintf("%s%s%s", clearView, clearHistory, moveToStart)
// 	for i := range e.h - 1 {
// 		if i < len(e.b.lines) {
// 			n := e.b.buildNumber(i+1, len(fmt.Sprint(len(e.b.lines))))
// 			data += fmt.Sprintf("%s %s\n\r", n, string(e.b.lines[i].data))
// 		} else {
// 			data += fmt.Sprintf("%s~\n\r", emtpyLineSpases)
// 		}
// 	}
// 	switch e.curMode {
// 	case insert:
// 		data += "-- INSERT --"
// 		data += "\033[5 q"
// 		x := e.b.cursor.ofset + initialOfset + len(emtpyLineSpases)
// 		y := e.b.cursor.line + cursorLineOfset
// 		data += fmt.Sprintf("\033[%d;%dH", y, x)
// 	case command:
// 		data += fmt.Sprintf(":%s\033[%d;%dH", e.curCommand, e.h, len(e.curCommand)+initialOfset-1) // 1 is a magic number, just get use to it
// 		data += "\033[0 q"
// 	case normal:
// 		data += fmt.Sprintf("%s%s", emtpyLineSpases, e.file)
// 		data += "\033[0 q"
// 		data += fmt.Sprintf("\033[%d;%dH", e.b.cursor.line+cursorLineOfset, e.b.cursor.ofset+initialOfset+len(emtpyLineSpases))
// 	}
// 	fmt.Print(data)
// }

func (ui *UI) DrawNew(e *Editor) {
	emtpyLineSpases := buildSpaces(len(fmt.Sprint(len(e.b.lines))))
	data := ""
	data += fmt.Sprintf("%s%s%s", clearView, clearHistory, moveToStart)
	for i := ui.upperBorder; i < ui.lowerBorder-1; i++ {
		if i < len(e.b.lines) {
			n := e.b.buildNumber(i+1, len(fmt.Sprint(len(e.b.lines))))
			data += fmt.Sprintf("%s %s\n\r", n, string(e.b.lines[i].data))
		} else {
			data += fmt.Sprintf("%s~\n\r", emtpyLineSpases)
		}
	}
	x := e.b.cursor.ofset + initialOfset + len(emtpyLineSpases)
	y := e.ui.curRow + cursorLineOfset

	switch e.curMode {
	case insert:
		data += "-- INSERT --"
		data += "\033[5 q"
		data += fmt.Sprintf("\033[%d;%dH", y, x)
	case command:
		data += fmt.Sprintf(":%s\033[%d;%dH", e.curCommand, e.h, len(e.curCommand)+initialOfset-1) // 1 is a magic number, just get use to it
		data += "\033[0 q"
	case normal:
		data += fmt.Sprintf("%s%s", emtpyLineSpases, e.file)
		data += "\033[0 q"
		data += fmt.Sprintf("\033[%d;%dH", y, x)
	}
	fmt.Print(data)
}
