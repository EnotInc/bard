package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/editor/TUI/render"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/buffers"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"
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
// Message - shows at the bottom of the screen, in 'lower bar'. Used to display some messages
// Hello - ascii art of bard
// visual - anchor of wisual row and offset. Used to calculate visual selection between this point and the cursor
// render - an instance of render. Used to buld line with ansi sybols
// Redraw - chan, wich used to redraw the whole editor when window size is changed
type TUI struct {
	Redraw  chan bool
	render  *render.Renderer
	visual  *visual
	Message string
	Hello   [][]rune
	CurOff  int
	H       int
	W       int
	XScroll int
	CurRow  int
	YScroll int
	Save    bool
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
	numStr := strconv.Itoa(rn)
	numLen := len(numStr)
	var num strings.Builder

	if maxOffset <= enums.InitialOffset {
		maxOffset = enums.InitialOffset
	}
	fmt.Fprint(&num, ascii.Reset, strings.Repeat(" ", maxOffset-numLen))
	//fmt.Fprint(&num, ascii.Bold)

	theme := config.GetTheme().General

	if curLine+1 == n {
		fmt.Fprint(&num, theme.CurrentLine, numStr)
	} else {
		fmt.Fprint(&num, theme.LineNumber, numStr)
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

func (ui *TUI) fillSpaceWith(ln int) string {
	amount := max(ui.W-ln, 0)
	return strings.Repeat(" ", amount)
}

func (ui *TUI) fillSpace() string {
	amount := max(screen.W(), 0)
	return strings.Repeat(" ", amount)
}

// Little func, that used to build lower bar
func (ui *TUI) BuildLowerBar(x int, y int, curdata string, message string, cmd string) string {
	theme := config.GetTheme().General
	var data strings.Builder
	pos := fmt.Sprintf(" %d-%d ", x, y)
	fmt.Fprintf(&data, "%s%s%s %s%s%s ", theme.BottomBar, pos, curdata, theme.Message, message, theme.BottomBar)

	ln := 0
	if cmd != "" {
		fmt.Fprintf(&data, "<%s>", cmd)
		ln += len(cmd) + 2
	}
	fmt.Fprintf(&data, "%s", ui.fillSpace())

	return services.VisibleSubString(data.String(), 0, screen.W()-1)
}

// Used when used is is command mode. It simply moves curos to the bottom of the scneed and at the end of the input command
func (ui *TUI) BuildCommandBar(curdata string) string {
	theme := config.GetTheme().General
	var data strings.Builder
	cmd := theme.Command + " :" + theme.BottomBar
	fmt.Fprintf(&data, "%s%s%s%s\033[%d;%dH%s", theme.BottomBar, cmd, curdata, ui.fillSpaceWith(len(curdata)+2), ui.H, len(curdata)+enums.InitialOffset, ascii.Reset)

	return data.String()
}

func (ui *TUI) BuildLine(str []rune, show bool, start, end int, i int, isCurrent bool, isFirst bool, Type buffers.BufferType) string {
	if Type == buffers.Other {
		clear := services.ReplaceTabs(str)
		shift := services.CursorShift(str)
		return services.VisibleSubString(string(clear), start, end+shift)
	}

	l := ui.render.Render(str, i, show, isCurrent, isFirst, ui.XScroll, Type)

	return services.VisibleSubString(l, start, end)
}

func (ui *TUI) ResetRender() {
	ui.render.Reset()
}

func (ui *TUI) Center(l []rune) string {
	offset := max((ui.W-len(l))/2, 0)
	tabs := strings.Repeat(" ", offset)
	return tabs + string(l)
}

func (ui *TUI) BuildTabs(tabs []string, curTab int, show bool) string {
	theme := config.GetTheme().General
	if len(tabs) == 1 {
		icon := ""
		if len(tabs[0]) != 0 {
			icon = services.GetFileIcon(tabs[0])
			icon = strings.TrimPrefix(icon, "  ")
		}
		return fmt.Sprintf("%s[%s%s%s]", theme.Tab, icon, tabs[0], theme.Tab)
	}

	var s strings.Builder
	for i, tab := range tabs {
		icon := ""
		if len(tab) != 0 {
			icon = services.GetFileIcon(tab)
			icon = strings.TrimPrefix(icon, "  ")
		}

		color := ascii.ResetFg
		if i == curTab {
			color = ascii.Color(theme.Tab)
		}
		if show {
			fmt.Fprintf(&s, "%s[%d|%s%s%s]", color, i+1, icon, color, tab)
		} else {
			fmt.Fprintf(&s, "%s[%d %s%s]", color, i+1, icon, color)
		}
		fmt.Fprint(&s, ascii.ResetFg)
	}
	return s.String()
}

func (ui *TUI) PurgeCache() {
	ui.render.PurgeCache()
	screen.SendCall(calls.PurgeCache)
}

func (ui *TUI) ToggleRender() {
	ui.render.ToggleRender()
}

func (ui *TUI) ResizeRender(w int) {
	ui.render.Resize(w)
}
