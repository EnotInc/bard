package editor

import (
	"fmt"
	"strconv"
)

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
	data += fmt.Sprintf("%s", e.curType)
	data += fmt.Sprintf("%s%s%s", clearView, clearHistory, moveToStart)

	s := strconv.Itoa(len(e.buffer.lines))
	numberOfset := len(s)

	for i := range e.h - 1 {
		if i < len(e.buffer.lines) {
			n := buildNumber(i+1, numberOfset)
			data += fmt.Sprintf("%s %s\r\n", n, e.buffer.lines[i])
		} else {
			data += "~\r\n"
		}
	}
	switch e.curMode {
	case insert:
		data += fmt.Sprintf("--%s--", e.curMode)
	case command:
		data += fmt.Sprintf(":%s", e.curCommand)
	}
	if e.curMode != command {
		data += fmt.Sprintf("\t\t\t%d-%d", e.buffer.curLine, e.buffer.curOfset)
	}
	data += fmt.Sprintf("\033[%d;%dH", e.buffer.curLine+initialLineShift, e.buffer.curOfset+numberOfset+initialCurOfset)
	fmt.Print(data)
}
