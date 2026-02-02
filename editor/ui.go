package editor

import (
	"Enot/Bard/render"
	"fmt"
	"strconv"
	"strings"
)

type color string

const (
	resetFg color = "\033[0m"
	redFg   color = "\033[31m"
	grayFg  color = "\033[90m"
	//yellowFg   color = "\033[93m"
	yellowFg   color = "\033[33m"
	cyanFg     color = "\033[36m"
	cursorBloc       = "\x1b[0 q"
	cursorLine       = "\x1b[5 q"
)

const (
	cursorLineOfset = 1
	initialOfset    = 3
	cursorDataOfset = 20
)

type UI struct {
	rln         bool
	upperBorder int
	lowerBorder int
	leftBorder  int
	rightBorder int
	curRow      int
	curOff      int
	render      *render.Renderer
}

func InitUI(h int, w int) *UI {
	r := render.InitReder(w, h)
	ui := &UI{
		rln:         false,
		lowerBorder: h,
		upperBorder: 0,
		leftBorder:  0,
		rightBorder: w - initialOfset,
		curRow:      0,
		curOff:      0,
		render:      r,
	}
	return ui
}

func colorise(data string, c color) string {
	return fmt.Sprintf("%s%s%s", c, data, resetFg)
}

func (b *Buffer) buildNumber(n int, maxOfset int, rln bool) string {
	rn := n
	if rln && rn != b.cursor.line+1 {
		rn = b.cursor.line - n + 1
		if rn < 0 {
			rn *= -1
		}
	}
	numStr := strconv.Itoa(rn)
	numLen := len(numStr)
	num := ""

	if maxOfset <= initialOfset {
		maxOfset = initialOfset
	}
	for range maxOfset - numLen {
		num += " "
	}
	num = fmt.Sprintf("%s%s", num, numStr)

	if b.cursor.line+1 == n {
		num = colorise(num, yellowFg)
	} else {
		num = colorise(num, grayFg)
	}

	return num
}

func buildSpaces(maxOfset int) string {
	space := ""
	if maxOfset <= initialOfset {
		maxOfset = initialOfset
	}
	for range maxOfset - 1 {
		space += " "
	}
	return space
}

func (e *Editor) buildLowerBar(curdata string) string {
	var data = ""
	data += fmt.Sprintf(" %d-%d ", e.b.cursor.line+1, e.b.cursor.ofset+1)
	data += fmt.Sprintf("%s %s%s%s", curdata, redFg, e.message, resetFg)
	return data
}

func (ui *UI) Draw(e *Editor) {
	emtpyLineSpases := buildSpaces(len(fmt.Sprint(len(e.b.lines))))
	maxNumLen := len(fmt.Sprint(len(e.b.lines)))

	var data strings.Builder
	fmt.Fprintf(&data, "%s%s%s", clearView, clearHistory, moveToStart)

	for i := ui.upperBorder; i < ui.lowerBorder-1; i++ {
		isCurLine := e.b.cursor.line == i

		start := e.ui.leftBorder
		end := e.ui.rightBorder - len(emtpyLineSpases)

		if i < len(e.b.lines) {

			str := e.b.lines[i].data
			if len(str) <= end {
				end = len(str)
			}
			if len(str) < start {
				start = 0
				end = 0
				str = []rune{}
			}

			n := e.b.buildNumber(i+1, maxNumLen, ui.rln)
			var l = ""
			if e.isMdFile {
				l = ui.render.RednerMarkdownLine(str[start:end], isCurLine)
			} else {
				l = string(str[start:end])
			}

			fmt.Fprintf(&data, "%s %s\n\r", n, l)
			//fmt.Fprintf(&data, "%s %s\n\r", n, string(str[start:end]))
		} else {
			fmt.Fprintf(&data, "%s\n\r", colorise("~", cyanFg))
		}
	}
	x := e.ui.curOff + initialOfset + len(emtpyLineSpases)
	y := e.ui.curRow + cursorLineOfset

	switch e.curMode {
	case insert:
		fmt.Fprintf(&data, "%s", e.buildLowerBar("-- INSERT --"))
		fmt.Fprintf(&data, cursorLine)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	case command:
		fmt.Fprintf(&data, "%s%s\033[%d;%dH", colorise(" :", yellowFg), e.curCommand, e.h, len(e.curCommand)+initialOfset)
		fmt.Fprintf(&data, cursorBloc)
	case normal:
		cursorPos := fmt.Sprintf("[%s]", e.file)
		fmt.Fprintf(&data, "%s", e.buildLowerBar(cursorPos))
		fmt.Fprintf(&data, cursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	case visual, visual_line:
		fmt.Fprintf(&data, "%s", e.buildLowerBar(string(e.curMode)))
		fmt.Fprintf(&data, cursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	}
	fmt.Print(data.String())
}
