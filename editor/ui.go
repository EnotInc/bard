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
	redFg      color = "\033[31m"
	grayFg     color = "\033[90m"
	yellowFg   color = "\033[33m"
	cyanFg     color = "\033[36m"
	cursorBloc       = "\x1b[0 q"
	cursorLine       = "\x1b[5 q"
)

const (
	cursorLineOfset = 1
	initialOfset    = 3
)

type UI struct {
	rln     bool
	xScroll int
	yScroll int
	curRow  int
	curOff  int
	w, h    int
	render  *render.Renderer
}

func InitUI(h int, w int) *UI {
	r := render.InitReder(w, h)
	ui := &UI{
		rln:     false,
		xScroll: 0,
		yScroll: 0,
		curRow:  0,
		curOff:  0,
		w:       w,
		h:       h,
		render:  r,
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
	num = strings.Repeat(" ", maxOfset-numLen)
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
	space = strings.Repeat(" ", maxOfset-1)
	return space
}

func (e *Editor) buildLowerBar(curdata string) string {
	var data = ""
	data += fmt.Sprintf(" %d-%d ", e.b.cursor.line+1, e.b.cursor.ofset+1)
	data += fmt.Sprintf("%s %s%s%s", curdata, redFg, e.message, resetFg)

	if e.subCmd != "" {
		data += fmt.Sprintf("<%s>", e.subCmd)
	}

	return data
}

func visibleSubString(text string, start int, end int) string {
	var res strings.Builder
	visibleCount := 0
	inEscape := false
	var escapeSeq strings.Builder

	for _, r := range text {
		if r == '\033' {
			inEscape = true
			escapeSeq.Reset()
			escapeSeq.WriteRune(r)
			continue
		}
		if inEscape {
			escapeSeq.WriteRune(r)
			if r == 'm' {
				inEscape = false
				if visibleCount >= start && visibleCount <= start+end {
					res.WriteString(escapeSeq.String())
				}
			}
			continue
		}
		if visibleCount >= start && visibleCount <= start+end {
			res.WriteRune(r)
		}
		visibleCount++
	}

	return res.String()
}

func (ui *UI) Draw(e *Editor) {
	emtpyLineSpases := buildSpaces(len(fmt.Sprint(len(e.b.lines))))
	maxNumLen := len(fmt.Sprint(len(e.b.lines)))

	var data strings.Builder
	fmt.Fprintf(&data, "%s%s%s", clearView, clearHistory, moveToStart)

	upperBorder := ui.yScroll
	lowerBorder := ui.yScroll + ui.h - 1

	for i := upperBorder; i < lowerBorder; i++ {

		if i < len(e.b.lines) {
			isCurLine := e.b.cursor.line == i

			start := ui.xScroll
			end := ui.w - initialOfset - len(emtpyLineSpases)

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
			var diff = 0
			if e.isMdFile {
				l, diff = ui.render.RednerMarkdownLine(str, isCurLine)
				if isCurLine {
					diff = 0
				}
				l = visibleSubString(l, start, end-diff)
				l += string(resetFg)
			} else {
				l = string(str[start:end])
			}

			fmt.Fprintf(&data, "%s %s\n\r", n, l)
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
		fmt.Fprintf(&data, "%s%s\033[%d;%dH", colorise(" :", yellowFg), e.command, ui.h, len(e.command)+initialOfset)
		fmt.Fprintf(&data, cursorBloc)
	case normal:
		cursorPos := fmt.Sprintf("[%s]", e.file)
		fmt.Fprintf(&data, "%s", e.buildLowerBar(cursorPos))
		fmt.Fprintf(&data, cursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	}
	fmt.Print(data.String())
}
