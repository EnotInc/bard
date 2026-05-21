package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/render"
	"golang.org/x/term"
)

type visual struct {
	offset int
	line   int
}

// so this is a struct where I store a data about visual intermretation of bard
// XScroll - stores a upped border of editor 'view window'
// YScroll - stores a left border of editor 'view window'
// CurRow - visual row where cursor is located
// CurOff - visual offset where cursor is located
// W, H - width and height of terminal window
// Save - is terminal save (is is big enough)
// ShowHello - used to how 'hello message' and bard logo in empty editor
// Message - shows at the bottom of the screen, in 'lower bar'. Used to display some messages
// Hello - ascii art of bard
// visual - anchor of wisual row and offset. Used to calculate visual selection between this point and the cursor
// render - an instance of render. Used to buld line with ansi sybols
// Redraw - chan, wich used to redraw the whole editor when window size is changed
type TUI struct {
	XScroll   int
	YScroll   int
	CurRow    int
	CurOff    int
	W, H      int
	Save      bool
	ShowHello bool
	Message   string
	Hello     [][]rune
	visual    *visual
	render    *render.Renderer
	theme     *config.General
	Redraw    chan bool
}

func InitTUI(h int, w int, theme *config.Theme) *TUI {
	r := render.InitRender(w, h, theme)
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
		theme:   &theme.General,
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
			tui.render.Resize(w)
		}
		tui.Redraw <- false
	}
}

func (tui *TUI) resize(w int, h int) {
	tui.W = w
	tui.H = h
	tui.Redraw <- true
}

func (tui *TUI) MakeDirty() {
	tui.render.MakeDirty()
}

// this func is used to build pretty line numbers (represented with '.'):
// |..8  // foo func
// |..9  func foo() {
// |.10      bar()
// |.11      baz()
// |.12  }
func (ui *TUI) BuildNumber(curLine int, n int, maxOffset int, rln bool) string {
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
	fmt.Fprint(&num, ascii.Reset, strings.Repeat(" ", maxOffset-numLen))

	if curLine+1 == n {
		fmt.Fprint(&num, ui.theme.CurrentLine, numStr)
	} else {
		fmt.Fprint(&num, ui.theme.LineNumber, numStr)
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

func (ui *TUI) fillSpaceWiht(ln int) string {
	amount := max(ui.W-ln, 0)
	return strings.Repeat(" ", amount)
}

func (ui *TUI) fillSpace() string {
	amount := max(ui.W, 0)
	return strings.Repeat(" ", amount)
}

// Little func, that used to build lower bar
func (ui *TUI) BuildLowerBar(x int, y int, curdata string, message string, cmd string) string {
	var data strings.Builder
	pos := fmt.Sprintf(" %d-%d ", x, y)
	fmt.Fprintf(&data, "%s%s%s %s%s%s ", ui.theme.BottomBar, pos, curdata, ui.theme.Message, message, ui.theme.BottomBar)

	ln := 0
	if cmd != "" {
		fmt.Fprintf(&data, "<%s>", cmd)
		ln += len(cmd) + 2
	}
	fmt.Fprintf(&data, "%s", ui.fillSpace())

	return VisibleSubString(data.String(), 0, ui.W-1)
}

// Used when used is is command mode. It simply moves curos to the bottom of the scneed and at the end of the input command
func (ui *TUI) BuildCommandBar(curdata string) string {
	var data strings.Builder
	cmd := ui.theme.Command + " :" + ui.theme.BottomBar
	fmt.Fprintf(&data, "%s%s%s%s\033[%d;%dH%s", ui.theme.BottomBar, cmd, curdata, ui.fillSpaceWiht(len(curdata)+2), ui.H, len(curdata)+enums.InitialOffset, ascii.Reset)

	return data.String()
}

// So here is where I build the actual line, including the ASCII escape sequences
// If I just use line.data[start:end], I'll get something like this:
// ```
// 033[0m and some text
// ```
// Here I just ignore the escape sequences and don't count them, so I can use them
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

func (ui *TUI) BuildLine(str []rune, show bool, start, end int, i int, isCurrent bool, isFirst bool, isRender bool) string {
	if !isRender { // returning stripped text if render is of (or it's not a md file)
		return VisibleSubString(string(str), start, end)
	}

	var l = ""
	// diff is used for calculating the size of the line, where markdown symbols are hidden
	var diff = 0
	l, diff = ui.render.Render(str, i, show, isCurrent, isFirst)
	if end < len(str) {
		diff = 0
	}
	l = VisibleSubString(l, start, end-diff)

	return l
}

func (ui *TUI) ResetRender() {
	ui.render.Reset()
}

func (ui *TUI) Center(l []rune) string {
	center := max((ui.W-len(l))/2, 0)
	tabs := strings.Repeat(" ", center)
	return tabs + string(l)
}

func (ui *TUI) BuildTabs(tabs []string, curTab int, show bool) string {
	if len(tabs) == 1 {
		return fmt.Sprintf("%s[%s]", ui.theme.Tab, tabs[0])
	}

	var s strings.Builder
	for i, tab := range tabs {
		if i == curTab {
			fmt.Fprint(&s, ui.theme.Tab)
		}
		if show {
			fmt.Fprintf(&s, "[%d|%s]", i+1, tab)
		} else {
			fmt.Fprintf(&s, "[%d]", i+1)
		}
		fmt.Fprint(&s, ascii.Reset, ui.theme.BottomBar)
	}
	return s.String()
}

func (ui *TUI) PurgeCache() {
	ui.render.PurgeCache()
}
