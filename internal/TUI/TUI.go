package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/render"
	"golang.org/x/term"
)

type visual struct {
	offset int
	line   int
}

type TUI struct {
	XScroll   int
	YScroll   int
	CurRow    int
	CurOff    int
	W, H      int
	Save      bool
	ShowHello bool // show ascii art in empty Bard
	Message   string
	Hello     [][]rune // ascii art in empty Bard
	visual    *visual
	render    *render.Renderer
	Redraw    chan bool
}

func InitTUI(h int, w int) *TUI {
	r := render.InitRender(w, h)
	v := &visual{line: 0, offset: 0}
	ui := &TUI{
		XScroll: 0,
		YScroll: 0,
		CurRow:  0,
		CurOff:  0,
		Save:    true,
		W:       w,
		H:       h,
		visual:  v,
		render:  r,
		Redraw:  make(chan bool, 1),
	}
	return ui
}

// This function is called in the main.go file in a goroutine.
// Here I just recalculate the terminal size and adjust Bard to it
func (tui *TUI) TermSizeMonitor(fdOut int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var last_w, last_h = tui.W, tui.H

	for range ticker.C {
		w, h, err := term.GetSize(fdOut)
		if err != nil {
			continue
		}

		if last_w != w || last_h != h {
			last_w = w
			last_h = h

			tui.resize(w, h)
		}
		tui.Redraw <- false
	}
}

func (tui *TUI) resize(w int, h int) {
	tui.W = w
	tui.H = h
	tui.Redraw <- true
}

func BuildNumber(curLine int, n int, maxOffset int, rln bool) string {
	rn := n
	if rln && rn != curLine+1 {
		rn = curLine - n + 1
		if rn < 0 {
			rn *= -1
		}
	}
	var numStr = strconv.Itoa(rn)
	numLen := len(numStr)
	var num strings.Builder

	if maxOffset <= enums.InitialOffset {
		maxOffset = enums.InitialOffset
	}
	fmt.Fprint(&num, strings.Repeat(" ", maxOffset-numLen))

	if curLine+1 == n {
		fmt.Fprint(&num, ascii.YellowFg, numStr)
	} else {
		fmt.Fprint(&num, ascii.GrayFg, numStr)
	}
	fmt.Fprint(&num, ascii.Reset, " ")

	return num.String()
}

func BuildSpaces(maxOffset int) string {
	if maxOffset <= enums.InitialOffset {
		maxOffset = enums.InitialOffset
	}
	return strings.Repeat(" ", maxOffset-1)
}

func (ui *TUI) BuildLowerBar(x int, y int, curdata string, message string, cmd string) string {
	var data strings.Builder
	fmt.Fprintf(&data, " %d-%d ", x, y)
	fmt.Fprintf(&data, "%s %s %s%s", curdata, ascii.RedFg, message, ascii.Reset)

	if cmd != "" {
		fmt.Fprintf(&data, "<%s>", cmd)
	}

	return data.String()
}

func (ui *TUI) BuildCommandBar(curdata string) string {
	var data strings.Builder
	cmd := ascii.YellowFg.Str() + " :" + ascii.Reset.Str()
	fmt.Fprintf(&data, "%s%s\033[%d;%dH%s", cmd, curdata, ui.H, len(curdata)+enums.InitialOffset, ascii.Reset)

	return data.String()
}

/*
 * So here is where I build the actual line, including the ASCII escape sequences
 * If I just use line.data[start:end], I'll get something like this:
 *
 * 033[0m and some text
 *
 * Here I just ignore the escape sequences and don't count them, so I can use them
 */
func VisibleSubString(text string, start int, end int) string {
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

func (ui *TUI) BuildLine(str []rune, show bool, start, end int, i int) string {
	var l = ""
	// diff is used for calculating the size of the line, where markdown symbols are hidden
	var diff = 0
	l, diff = ui.render.Render(str, i, show)
	if show || end < len(str) {
		diff = 0
	}
	l = VisibleSubString(l, start, end-diff)

	return l
}

func (ui *TUI) ResetRender() {
	ui.render.Reset()
}

func (ui *TUI) Center(l []rune) string {
	center := (ui.W - len(l)) / 2
	if center < 0 {
		center = 0
	}
	tabs := strings.Repeat(" ", center)
	return tabs + string(l)
}

func (ui *TUI) BuildTabs(tabs []string, curTab int, show bool) string {
	var s strings.Builder
	for i, tab := range tabs {
		if i == curTab {
			fmt.Fprint(&s, ascii.Tab, ascii.Reset)
		}
		if show {
			fmt.Fprintf(&s, "[%d|%s]", i+1, tab)
		} else {
			fmt.Fprintf(&s, "[%d]", i+1)
		}
	}
	return s.String()
}
