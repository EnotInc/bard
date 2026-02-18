package editor

import (
	"Enot/Bard/render"
	"fmt"
	"strconv"
	"strings"
)

type color string

const (
	reset      color = "\033[0m"
	redFg      color = "\033[31m"
	grayFg     color = "\033[90m"
	yellowFg   color = "\033[33m"
	cyanFg     color = "\033[36m"
	startSel   color = "\033[100m"
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
	return fmt.Sprintf("%s%s", c, data /*, reset*/)
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
	num += string(reset)

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
	data += fmt.Sprintf("%s %s%s%s", curdata, redFg, e.message, reset)

	if e.subCmd != "" {
		data += fmt.Sprintf("<%s>", e.subCmd)
	}

	return data
}

/*
 * So here is where I build the actual line, including the ASCII escape sequences
 * If I just use line.data[start:end], I'll get something like this:
 *
 * 033[0m and some text
 *
 * Here I just ignore the escape sequences and don't count them, so I can use them
 */
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

// This function is used to add visual highlight to the selected lines
func (e *Editor) addVisual(l string, i int) string {
	var line string

	switch e.curMode {
	case visual:
		startOfset := e.b.visual.ofset
		startLine := e.b.visual.line

		endOfset := e.b.cursor.ofset
		endLine := e.b.cursor.line

		if startLine > endLine || (startLine == endLine && startOfset > endOfset) {
			startLine, endLine = endLine, startLine
			startOfset, endOfset = endOfset, startOfset
		}

		if len(e.b.lines[endLine].data) > 0 {
			endOfset++
		}

		if startLine == i && i == endLine {
			line = l[:startOfset] + string(startSel) + l[startOfset:endOfset] + string(reset) + l[endOfset:]
		} else if startLine < i && i < endLine {
			line = string(startSel) + l + string(reset)
		} else if startLine == i {
			line = l[:startOfset] + string(startSel) + l[startOfset:] + string(reset)
		} else if endLine == i {
			line = string(startSel) + l[:endOfset] + string(reset) + l[endOfset:]
		} else {
			line = l
		}

	case visual_line:
		startLine := e.b.visual.line
		endLine := e.b.cursor.line

		if startLine > endLine {
			startLine, endLine = endLine, startLine
		}

		line = string(startSel) + l + string(reset)
	}

	return line
}

func (ui *UI) buildLine(str []rune, show bool, start, end int) string {
	var l = ""
	// diff is used for calculating the size of the line, where markdown symbols are hidden
	var diff = 0
	l, diff = ui.render.RednerMarkdownLine(str, show)
	if show {
		diff = 0
	}
	l = visibleSubString(l, start, end-diff)

	return l
}

func (ui *UI) Draw(e *Editor) {
	emtpyLineSpases := buildSpaces(len(fmt.Sprint(len(e.b.lines))))
	maxNumLen := len(fmt.Sprint(len(e.b.lines)))

	// data - is one long string that turns into the TUI
	var data strings.Builder

	// Clearing the terminal
	fmt.Fprintf(&data, "%s%s%s", clearView, clearHistory, moveToStart)

	upperBorder := ui.yScroll
	lowerBorder := ui.yScroll + ui.h - 1

	isVisual := e.curMode == visual

	// Working only with visible lines
	for i := upperBorder; i < lowerBorder; i++ {
		if i < len(e.b.lines) {
			show := e.b.cursor.line == i || isVisual

			// This 2 variables is used to get the horizontal borders of visible content
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
			if e.isMdFile {
				switch e.curMode {
				case visual, visual_line:
					// This if statement lets me render both selected lines with highlights, and not selected with markdown render
					if (i >= e.b.visual.line && i <= e.b.cursor.line) || (i <= e.b.visual.line && i >= e.b.cursor.line) {
						l = e.addVisual(string(str[start:end]), i)
					} else {
						l = ui.buildLine(str, show, start, end)
					}
				// Some other modes can use different logic for rendering, but now I just call the default for non-visual or visual_line modes
				default:
					l = ui.buildLine(str, show, start, end)
				}

				l += string(reset)
			} else {
				if e.curMode == visual || e.curMode == visual_line {
					l = e.addVisual(string(str[start:end]), i)
				} else {
					l = string(str[start:end])
				}
			}

			// Here is where I add the line to the main data string
			fmt.Fprintf(&data, "%s %s\n\r", n, l)
		} else {
			// If the line is empty, I just add the '~' symbol
			fmt.Fprintf(&data, "%s\n\r", colorise("~", cyanFg))
		}
	}

	// Calculation the visual position of cursor
	x := e.ui.curOff + initialOfset + len(emtpyLineSpases)
	y := e.ui.curRow + cursorLineOfset

	//Different modes have different information on the last line
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
	case visual, visual_line:
		fmt.Fprintf(&data, "%s", e.buildLowerBar(fmt.Sprintf("-- %s --", e.curMode)))
		fmt.Fprintf(&data, cursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	}

	// And at the end - print the data
	fmt.Print(data.String())
}
