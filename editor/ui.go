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
	initialOfset = 2
)

func colorise(data string, color string) string {
	return fmt.Sprintf("%s%s", color, data)
}

func buildNumber(n int, maxOfset int) string {
	numStr := strconv.Itoa(n)
	numLen := len(numStr)
	num := ""
	for range maxOfset - numLen {
		num += " "
	}
	num += numStr
	return fmt.Sprintf("%s", num)
}

func Draw(e *Editor) {
	data := ""
	data += fmt.Sprintf("%s%s%s", clearView, clearHistory, moveToStart)
	for i := range e.h - 1 {
		if i < len(e.b.lines) {
			n := buildNumber(i+1, len(fmt.Sprint(len(e.b.lines))))
			data += fmt.Sprintf("%s %s\n\r", n, string(e.b.lines[i].data))
		} else {
			data += "~\n\r"
		}
	}
	switch e.curMode {
	case insert:
		data += "-- INSERT --"
		data += fmt.Sprintf("\033[%d;%dH", e.b.cursor.line, e.b.cursor.ofset+initialCurOfset)
	case command:
		data += fmt.Sprintf(":%s\033[%d;%dH", e.curCommand, e.h, len(e.curCommand)+initialOfset)
	case normal:
		data += fmt.Sprintf("\033[%d;%dH", e.b.cursor.line, e.b.cursor.ofset+initialOfset)
	}
	fmt.Print(data)
}
