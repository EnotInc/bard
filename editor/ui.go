package editor

import (
	"fmt"
	"strconv"
)

type color string

const (
	red    color = ""
	green  color = ""
	gray   color = ""
	yellow color = ""
	black  color = ""
)

const (
	initialOfset    = 3
	cursorLineOfset = 1
)

func colorise(data string, color string) string {
	return fmt.Sprintf("%s%s", color, data)
}

func buildNumber(n int, maxOfset int) string {
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

func Draw(e *Editor) {
	emtpyLineSpases := buildSpaces(len(fmt.Sprint(len(e.b.lines))))
	data := ""
	data += fmt.Sprintf("%s%s%s", clearView, clearHistory, moveToStart)
	for i := range e.h - 1 {
		if i < len(e.b.lines) {
			n := buildNumber(i+1, len(fmt.Sprint(len(e.b.lines))))
			data += fmt.Sprintf("%s %s\n\r", n, string(e.b.lines[i].data))
		} else {
			data += fmt.Sprintf("%s~\n\r", emtpyLineSpases)
		}
	}
	switch e.curMode {
	case insert:
		data += "-- INSERT --"
		data += "\033[5 q"
		x := e.b.cursor.ofset + initialOfset + len(emtpyLineSpases)
		y := e.b.cursor.line + cursorLineOfset
		data += fmt.Sprintf("\033[%d;%dH", y, x)
	case command:
		data += fmt.Sprintf(":%s\033[%d;%dH", e.curCommand, e.h, len(e.curCommand)+initialOfset-1) // 1 is a magic number, just get use to it
		data += "\033[0 q"
	case normal:
		data += fmt.Sprintf("%s%s", emtpyLineSpases, e.file)
		data += "\033[0 q"
		data += fmt.Sprintf("\033[%d;%dH", e.b.cursor.line+cursorLineOfset, e.b.cursor.ofset+initialOfset+len(emtpyLineSpases))
	}
	fmt.Print(data)
}
