package editor

import (
	"Enot/Bard/render"
	"fmt"
	"strconv"
	"strings"
)

type color string

const (
	resetFg    color = "\033[0m"
	yellowFg   color = "\033[93m"
	cursorBloc       = "\x1b[0 q"
	cursorLine       = "\x1b[5 q"
)

const (
	initialOfset    = 3
	initialCurOfset = 1
	cursorLineOfset = 1
	cursorDataOfset = 20
)

type UI struct {
	rln         bool
	upperBorder int
	lowerBorder int
	curRow      int
	render      *render.Renderer
}

func InitUI(h int, w int) *UI {
	r := render.InitReder(w)
	ui := &UI{
		rln:         false,
		lowerBorder: h,
		upperBorder: 0,
		curRow:      0,
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
		num = colorise(num, resetFg)
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
	freeSpace := e.w - cursorDataOfset - len(curdata)
	var data = curdata
	for range freeSpace {
		data += " "
	}
	data += fmt.Sprintf("%d-%d", e.b.cursor.line+1, e.b.cursor.ofset+1)
	return data
}

func (ui *UI) Draw(e *Editor) {
	emtpyLineSpases := buildSpaces(len(fmt.Sprint(len(e.b.lines))))
	maxNumLen := len(fmt.Sprint(len(e.b.lines)))

	var data strings.Builder
	fmt.Fprintf(&data, "%s%s%s", clearView, clearHistory, moveToStart)

	for i := ui.upperBorder; i < ui.lowerBorder-1; i++ {
		isCurLine := e.b.cursor.line == i

		if i < len(e.b.lines) {
			n := e.b.buildNumber(i+1, maxNumLen, ui.rln)
			l := ""
			switch e.curMode {
			case visual, visual_line:
				l = ui.render.RenderVisualLine(e.b.lines[i].data, i, isCurLine)
			default:
				l = ui.render.RednerMarkDownLine(e.b.lines[i].data, isCurLine)
			}
			fmt.Fprintf(&data, "%s %s\n\r", n, l)
		} else {
			fmt.Fprintf(&data, "%s~\n\r", emtpyLineSpases)
		}
	}
	x := e.b.cursor.ofset + initialOfset + len(emtpyLineSpases) + initialCurOfset - 1
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
		cursorPos := fmt.Sprintf("%s[%s]", emtpyLineSpases, e.file)
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
