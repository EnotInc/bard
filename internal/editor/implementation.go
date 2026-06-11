package editor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/EnotInc/Bard/config"
	tui "github.com/EnotInc/Bard/internal/editor/TUI"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/screen"
	"github.com/EnotInc/Bard/internal/services"

	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

func (e *Editor) DrawStatusBar(withBorder bool) string {
	var data strings.Builder

	// NOTE: this could be borke if editor will be above another tile. But for now it's working
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
		fmt.Fprint(&data, e.tui.BuildCommandBar(string(e.cmd.command)))
		fmt.Fprintf(&data, ascii.CursorBloc)

	case mode.Normal:
		tabNames := ""

		if !withBorder {
			var tabs []string
			for _, t := range e.b {
				tab := t.Title
				if t.Title != "" {
					tab = filepath.Base(t.Title)
				}
				tabs = append(tabs, tab)
			}
			cfg := config.GetConfig()
			tabNames = e.tui.BuildTabs(tabs, e.curBuffer, cfg.TabNames)
		}

		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(posx, posy, tabNames, e.tui.Message, e.subCmd))
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
	maxNumLen := len(fmt.Sprint(len(e.b[e.curBuffer].Lines)))
	l, _ := e.drawRenderedLine(index+upperBorder, upperBorder, maxNumLen)

	return l
}

func (e *Editor) drawRenderedLine(i int, upperBorder int, maxNumLen int) (string, bool) {
	cfg := config.GetConfig()
	buf := e.b[e.curBuffer]
	show := buf.Cursor.Line() == i || cfg.ShowMD
	isFirst := i == upperBorder

	var l strings.Builder
	var keep bool

	if i < len(buf.Lines) { // rendering line
		var content strings.Builder

		start := e.tui.XScroll
		end := e.tui.W - enums.InitialOffset - len(e.emtpyLineSpases)

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

func (e *Editor) Handle(key rune) {
	switch e.curMode {
	case mode.Normal:
		if e.IsGeneralMove(key) {
			e.GeneralCase(key)
		} else {
			e.caseNormal(key)
		}
	case mode.Visual:
		if e.IsGeneralMove(key) {
			e.GeneralCase(key)
		} else {
			e.caseVisual(key)
		}
	case mode.Visual_line:
		if e.IsGeneralMove(key) {
			e.GeneralCase(key)
		} else {
			e.caseVisualLine(key)
		}
	case mode.Command:
		e.caseCommand(key)
	case mode.Insert:
		e.caseInsert(key)
	case mode.Replace:
		e.caseReplaceMode(key)
	default:
		screen.Exit(1)
	}

	e.setUiCursor()
}

func (e *Editor) GetCursor(withBorder bool) (int, int) {
	var x int
	var y int

	if e.curMode == mode.Command {
		x = len(e.cmd.command) + enums.InitialOffset - 1
		y = e.tui.H

	} else {
		x = e.tui.CurOff + enums.InitialOffset + len(e.emtpyLineSpases)
		y = e.tui.CurRow + enums.CursorOffset
	}

	if !withBorder {
		x += 1
	}

	return x, y
}

func (e *Editor) SetTitle() string {
	var tabs []string
	for _, t := range e.b {

		tab := t.Title
		if t.Title != "" {
			tab = filepath.Base(t.Title)
		}

		tabs = append(tabs, tab)
	}
	cfg := config.GetConfig()
	return e.tui.BuildTabs(tabs, e.curBuffer, cfg.TabNames)
}

func (e *Editor) PreDraw() {
	e.setUiCursor()
	e.emtpyLineSpases = tui.BuildSpaces(len(fmt.Sprint(len(e.b[e.curBuffer].Lines))))
	for i := range e.tui.YScroll {
		curLine := string(e.b[e.curBuffer].Lines[i].Data)
		if strings.HasPrefix(curLine, "```") {
			e.tui.ToggleRender()
		}
	}
}

func (e *Editor) Resize(w, h int) {
	e.tui.PurgeCache()
	e.tui.ResizeRender(w)
	e.tui.W = w
	e.tui.H = h
}
