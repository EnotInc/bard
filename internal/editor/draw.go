package editor

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
	tui "github.com/EnotInc/Bard/internal/editor/TUI"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"

	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

func (e *Editor) DrawStatusBar(withBorder bool) string {
	var data strings.Builder

	// NOTE: this coulb be borke if editor will be above another tile. But for now it's working
	fmt.Fprintf(&data, "\033[%d;1H", e.tui.H+1)

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

	case mode.Replace:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorUnderline)

	case mode.Command:
		//FIXME: figure out how to draw a ghost cursor
		fmt.Fprint(&data, e.tui.BuildCommandBar(string(e.cmd.command)))
		fmt.Fprintf(&data, ascii.CursorBloc)

	case mode.Normal:
		tab := ""

		if !withBorder {
			var tabs []string
			for _, t := range e.b {
				tabs = append(tabs, t.Title)
			}
			cfg := config.GetConfig()
			tab = e.tui.BuildTabs(tabs, e.curBuffer, cfg.TabNames)
		}

		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, tab, e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)

	case mode.Visual, mode.Visual_line:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
	}

	e.tui.ResetRender()
	fmt.Fprint(&data, ascii.Reset, ascii.ShowCursor)

	e.tui.Message = ""

	return data.String()
}

func (e *Editor) DrawLineAt(index int) string {
	upperBorder := e.tui.YScroll
	// TODO: move this shit out, this could be calculated at the start of draw call
	emtpyLineSpases := tui.BuildSpaces(len(fmt.Sprint(len(e.b[e.curBuffer].Lines))))
	maxNumLen := len(fmt.Sprint(len(e.b[e.curBuffer].Lines)))
	l, _ := e.drawRenderedLine(index+upperBorder, upperBorder, emtpyLineSpases, maxNumLen)

	return l
}

func (e *Editor) drawRenderedLine(i int, upperBorder int, emtpyLineSpases string, maxNumLen int) (string, bool) {
	cfg := config.GetConfig()
	buf := e.b[e.curBuffer]
	show := buf.Cursor.Line() == i || cfg.ShowMD
	isFirst := i == upperBorder

	var l strings.Builder
	var keep bool

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
