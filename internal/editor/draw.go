package editor

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/hash"

	tui "github.com/EnotInc/Bard/internal/TUI"
)

func (e *Editor) DrawDiff() {
	emtpyLineSpases := tui.BuildSpaces(len(fmt.Sprint(len(e.b[e.curBuffer].Lines))))
	maxNumLen := len(fmt.Sprint(len(e.b[e.curBuffer].Lines)))

	var diff strings.Builder

	fmt.Fprint(&diff, ascii.HideCursor, ascii.MoveToStart)

	upperBorder := e.tui.YScroll
	lowerBorder := e.tui.YScroll + e.tui.H - 1

	for i := upperBorder; i < lowerBorder; i++ {
		l := e.getRenderedLine(i, upperBorder, emtpyLineSpases, maxNumLen)
		curHash := hash.GetHash(l)
		if oldHash, ok := e.hash[i-upperBorder]; !ok || (ok && curHash != oldHash) {
			fmt.Fprintf(&diff, "\033[%d;1H\033[0K", i-upperBorder+1)
			fmt.Fprint(&diff, l)
			e.hash[i-upperBorder] = curHash
		}
	}

	status := e.getStatusBar(emtpyLineSpases, lowerBorder-upperBorder)
	fmt.Fprint(&diff, status)
	fmt.Print(diff.String())
}

func (e *Editor) getStatusBar(emtpyLineSpases string, lastLine int) string {
	var data strings.Builder

	fmt.Fprintf(&data, "\033[%d;1H", lastLine+1)

	x := e.tui.CurOff + enums.InitialOffset + len(emtpyLineSpases)
	y := e.tui.CurRow + enums.CursorOffset

	cursor := e.b[e.curBuffer].Cursor
	posx := cursor.Offset() + enums.CursorOffset
	posy := cursor.Line() + enums.CursorOffset

	fmt.Fprintf(&data, "%s", ascii.Reset)

	if e.b[e.curBuffer].IsReadOnly && e.tui.Message == "" {
		e.tui.Message = "read only file"
	}

	// Different modes have different information on the last line
	switch e.curMode {
	case enums.Insert:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorLine)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case enums.Replace:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorUnderline)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case enums.Command:
		fmt.Fprint(&data, e.tui.BuildCommandBar(string(e.command)))
		fmt.Fprintf(&data, ascii.CursorBloc)

	case enums.Normal:
		var tabs []string
		for _, t := range e.b {
			tabs = append(tabs, t.Title)
		}
		cursorPos := e.tui.BuildTabs(tabs, e.curBuffer, e.c.TabNames)
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, cursorPos, e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case enums.Visual, enums.Visual_line:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	}

	e.tui.ResetRender()
	fmt.Fprint(&data, ascii.Reset, ascii.ShowCursor)

	return data.String()
}

func (e *Editor) getRenderedLine(i int, upperBorder int, emtpyLineSpases string, maxNumLen int) string {
	buf := e.b[e.curBuffer]
	show := buf.Cursor.Line() == i || e.c.ShowMD
	isFirst := i == upperBorder

	var l strings.Builder

	// This 2 variables are used to get the horizontal borders of the visible content
	if i < len(buf.Lines) { // rendering line
		var content strings.Builder

		start := e.tui.XScroll
		end := e.tui.W - enums.InitialOffset - len(emtpyLineSpases)

		str := buf.Lines[i].Data
		if len(str) <= end {
			end = len(str)
		}
		if len(str) < start {
			start = 0
			end = 0
			str = []rune{}
		}

		n := e.tui.BuildNumber(buf.Cursor.Line(), i+1, maxNumLen, e.c.RLN)

		isRender := e.b[e.curBuffer].IsMdFile && e.c.Render
		switch e.curMode {
		case enums.Visual, enums.Visual_line:
			// This `if statement` let me render both selected lines with highlights, and not selected with markdown render
			if (i >= buf.Visual.Line() && i <= buf.Cursor.Line()) || (i <= e.b[e.curBuffer].Visual.Line() && i >= e.b[e.curBuffer].Cursor.Line()) {
				visual := e.tui.AddVisual(e.curMode, str, i, buf.Visual.Offset(), buf.Visual.Line(), buf.Cursor.Offset(), buf.Cursor.Line(), len(buf.Lines[buf.Cursor.Line()].Data), isRender)
				fmt.Fprint(&content, tui.VisibleSubString(visual, start, end))
			} else {
				fmt.Fprint(&content, e.tui.BuildLine(str, show, start, end, i, i == buf.Cursor.Line(), isFirst, isRender))
			}
		// Some other modes can use different logic for rendering, but now I just call the default for non-visual or visual_line modes
		default:
			fmt.Fprint(&content, e.tui.BuildLine(str, show, start, end, i, i == buf.Cursor.Line(), isFirst, isRender))
		}

		// Here is where I add the line to the main data string
		fmt.Fprint(&content, ascii.Reset)
		fmt.Fprint(&l, n, content.String())
	} else { // getting empty line
		if e.tui.ShowHello {
			fmt.Fprint(&l, ascii.Reset, e.theme.General.EmptyLine, "~", ascii.Reset, e.tui.Center(e.tui.GetASCIIInfo(i)))
		} else {
			fmt.Fprint(&l, ascii.Reset, e.theme.General.EmptyLine, "~")
		}
	}

	return l.String()
}
