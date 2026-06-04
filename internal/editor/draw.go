package editor

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/services"

	tui "github.com/EnotInc/Bard/internal/TUI"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

func (e *Editor) DrawDiff() {
	emtpyLineSpases := tui.BuildSpaces(len(fmt.Sprint(len(e.b[e.curBuffer].Lines))))
	maxNumLen := len(fmt.Sprint(len(e.b[e.curBuffer].Lines)))

	var diff strings.Builder

	fmt.Fprint(&diff, ascii.HideCursor, ascii.MoveToStart)

	upperBorder := e.tui.YScroll
	lowerBorder := e.tui.YScroll + e.tui.H - 1

	i := upperBorder
	for j := range lowerBorder {
		if j < upperBorder {
			curLine := string(e.b[e.curBuffer].Lines[j].Data)
			if strings.HasPrefix(curLine, "```") {
				e.tui.ToggleRender()
			}
			continue
		}

		l, keep := e.drawRenderedLine(i, upperBorder, emtpyLineSpases, maxNumLen)
		curHash := services.GetHash(l)
		oldHash, ok := e.hash[i-upperBorder]

		// if row is 1 of cursor position, or hash isn't the same as prev render, or hash wasn't calculated - draw line
		if keep || !ok || (ok && curHash != oldHash) || (i-upperBorder == e.tui.CurRow) {
			fmt.Fprintf(&diff, "\033[%d;1H\033[0K", i-upperBorder+1)
			fmt.Fprint(&diff, l)
			e.hash[i-upperBorder] = curHash
		}
		i++
	}

	status := e.drawStatusBar(emtpyLineSpases, lowerBorder-upperBorder)
	fmt.Fprint(&diff, status)
	fmt.Print(diff.String())
}

func (e *Editor) drawStatusBar(emtpyLineSpases string, lastLine int) string {
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
	case mode.Insert:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorLine)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case mode.Replace:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorUnderline)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case mode.Command:
		fmt.Fprintf(&data, "\033[%d;%dH%s\033[%d;1H", y, x, ascii.Cursor, lastLine+1) // adding cursor as unicode symbol on last visual position, and returning it to the last line
		fmt.Fprint(&data, e.tui.BuildCommandBar(string(e.cmd.command)))
		fmt.Fprintf(&data, ascii.CursorBloc)

	case mode.Normal:
		var tabs []string
		for _, t := range e.b {
			tabs = append(tabs, t.Title)
		}
		cfg := config.GetConfig()
		cursorPos := e.tui.BuildTabs(tabs, e.curBuffer, cfg.TabNames)
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, cursorPos, e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case mode.Visual, mode.Visual_line:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	}

	e.tui.ResetRender()
	fmt.Fprint(&data, ascii.Reset, ascii.ShowCursor)

	return data.String()
}

func (e *Editor) drawRenderedLine(i int, upperBorder int, emtpyLineSpases string, maxNumLen int) (string, bool) {
	cfg := config.GetConfig()
	buf := e.b[e.curBuffer]
	show := buf.Cursor.Line() == i || cfg.ShowMD
	isFirst := i == upperBorder

	var l strings.Builder
	var keep bool

	// This 2 variables are used to get the horizontal borders of the visible content
	if i < len(buf.Lines) { // rendering line
		var content strings.Builder

		start := e.tui.XScroll
		end := e.tui.W - enums.InitialOffset - len(emtpyLineSpases)

		str := buf.Lines[i].Data

		n := e.tui.BuildNumber(buf.Cursor.Line(), i+1, maxNumLen, cfg.RLN)

		isRender := e.b[e.curBuffer].IsMdFile && cfg.Render

		var data string

		switch e.curMode {
		case mode.Visual, mode.Visual_line:
			// This `if statement` let me render both selected lines with highlights, and not selected with markdown render
			if (i >= buf.Visual.Line() && i <= buf.Cursor.Line()) || (i <= e.b[e.curBuffer].Visual.Line() && i >= e.b[e.curBuffer].Cursor.Line()) {
				visual := e.tui.AddVisual(e.curMode,
					str, i,
					buf.Visual.Offset(),
					buf.Visual.Line(),
					buf.Cursor.Offset(),
					buf.Cursor.Line(),
					isRender)

				fmt.Fprint(&content, services.VisibleSubString(visual, start, end))
			} else {
				data, keep = e.tui.BuildLine(str, show, start, end, i, i == buf.Cursor.Line(), isFirst, isRender)
				fmt.Fprint(&content, data)
			}
		// Some other modes can use different logic for rendering, but now I just call the default for non-visual or visual_line modes
		default:
			data, keep = e.tui.BuildLine(str, show, start, end, i, i == buf.Cursor.Line(), isFirst, isRender)
			fmt.Fprint(&content, data)
		}

		// Here is where I add the line to the main data string
		fmt.Fprint(&content, ascii.Reset)
		fmt.Fprint(&l, n, content.String())
	} else { // getting empty line
		theme := config.GetTheme().General
		if e.tui.ShowHello {
			fmt.Fprint(&l, ascii.Reset, theme.EmptyLine, "~", ascii.Reset, e.tui.Center(e.tui.GetASCIIInfo(i)))
		} else {
			fmt.Fprint(&l, ascii.Reset, theme.EmptyLine, "~")
		}
	}

	return l.String(), keep
}

func (e *Editor) PurgeCache() {
	e.hash = make(map[int]uint32)
}
