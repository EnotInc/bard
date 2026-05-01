package editor

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/mode"

	tui "github.com/EnotInc/Bard/internal/TUI"
)

// About |Draw()|
// Main func to build tui and display it
// Line by line this function collects data from buffer, render raw text into markdown, accumulates it and prints it
// by collection all lines into one variable, I can avoid cursor blinking
// Curcor position is changed wish ascii escape sequence, and it calculates every time when this function is called
func (e *Editor) Draw() {
	emtpyLineSpases := tui.BuildSpaces(len(fmt.Sprint(len(e.b[e.curBuffer].Lines))))
	maxNumLen := len(fmt.Sprint(len(e.b[e.curBuffer].Lines)))

	buf := e.b[e.curBuffer]

	// data - is one long string that turns into the TUI
	var data strings.Builder

	// Clearing the terminal
	fmt.Fprint(&data, ascii.ClearView, ascii.ClearHistory, ascii.MoveToStart)

	upperBorder := e.tui.YScroll
	lowerBorder := e.tui.YScroll + e.tui.H - 1

	if e.IsChanged {
		e.tui.MakeDirty()
		e.IsChanged = false
	}

	// Working only with visible lines
	for i := upperBorder; i < lowerBorder; i++ {
		if i < len(buf.Lines) {
			show := buf.Cursor.Line() == i || e.c.Editor.ShowMD
			isFirst := i == upperBorder

			// This 2 variables are used to get the horizontal borders of the visible content
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

			n := e.tui.BuildNumber(buf.Cursor.Line(), i+1, maxNumLen, e.c.Editor.RLN)
			var l strings.Builder
			if e.b[e.curBuffer].IsMdFile && e.c.Editor.Render {
				switch e.curMode {
				case mode.Visual, mode.Visual_line:
					// This `if statement` let me render both selected lines with highlights, and not selected with markdown render
					if (i >= buf.Visual.Line() && i <= buf.Cursor.Line()) || (i <= e.b[e.curBuffer].Visual.Line() && i >= e.b[e.curBuffer].Cursor.Line()) {
						visual := e.tui.AddVisual(e.curMode, str, i, buf.Visual.Offset(), buf.Visual.Line(), buf.Cursor.Offset(), buf.Cursor.Line(), len(buf.Lines[buf.Cursor.Line()].Data))
						fmt.Fprint(&l, tui.VisibleSubString(visual, start, end))
					} else {
						fmt.Fprint(&l, e.tui.BuildLine(str, show, start, end, i, i == buf.Cursor.Line(), isFirst))
					}
				// Some other modes can use different logic for rendering, but now I just call the default for non-visual or visual_line modes
				default:
					fmt.Fprint(&l, e.tui.BuildLine(str, show, start, end, i, i == buf.Cursor.Line(), isFirst))
				}

				fmt.Fprint(&l, ascii.Reset.Str())
			} else {
				if e.curMode == mode.Visual || e.curMode == mode.Visual_line {
					visual := e.tui.AddVisual(e.curMode, str, i, buf.Visual.Offset(), buf.Visual.Line(), buf.Cursor.Offset(), buf.Cursor.Line(), len(buf.Lines[buf.Cursor.Line()].Data))
					fmt.Fprint(&l, tui.VisibleSubString(visual, start, end))
				} else {
					fmt.Fprint(&l, tui.VisibleSubString(string(str), start, end))
				}
			}

			// Here is where I add the line to the main data string
			fmt.Fprint(&data, n, l.String(), "\n\r")
			l.Reset()
		} else {
			// If the line is empty, I just add the '~' symbol
			if e.tui.ShowHello {
				fmt.Fprint(&data, e.c.Theme.General.EmptyLine, "~", ascii.Reset, e.tui.Center(e.tui.GetASCIIInfo(i)), "\n\r")
			} else {
				fmt.Fprint(&data, e.c.Theme.General.EmptyLine, "~", "\n\r")
			}
		}
	}

	// Calculating the visual position of the cursor
	x := e.tui.CurOff + enums.InitialOffset + len(emtpyLineSpases)
	y := e.tui.CurRow + enums.CursorLineOffset

	fmt.Fprintf(&data, "%s", ascii.Reset)

	if e.b[e.curBuffer].IsReadOnly && e.tui.Message == "" {
		e.tui.Message = "read only file"
	}

	// Different modes have different information on the last line
	switch e.curMode {
	case mode.Insert:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(x, y, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorLine)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case mode.Replace:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(x, y, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorUnderline)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case mode.Command:
		fmt.Fprint(&data, e.tui.BuildCommandBar(string(e.command)))
		fmt.Fprintf(&data, ascii.CursorBloc)

	case mode.Normal:
		var tabs []string
		for _, t := range e.b {
			tabs = append(tabs, t.Title)
		}
		cursorPos := e.tui.BuildTabs(tabs, e.curBuffer, e.c.Editor.TabNames)
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(x, y, cursorPos, e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case mode.Visual, mode.Visual_line:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(x, y, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	}

	e.tui.ResetRender()
	fmt.Fprintf(&data, "%s", ascii.Reset)

	// And at the end - print the data
	fmt.Print(data.String())
	data.Reset()
}
