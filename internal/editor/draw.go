package editor

import (
	tui "Enot/Bard/internal/TUI"
	"Enot/Bard/internal/ascii"
	"Enot/Bard/internal/enums"
	"Enot/Bard/internal/mode"
	"fmt"
	"strings"
)

func (e *Editor) Draw() {
	emtpyLineSpases := tui.BuildSpaces(len(fmt.Sprint(len(e.b[e.curBuffer].Lines))))
	maxNumLen := len(fmt.Sprint(len(e.b[e.curBuffer].Lines)))

	// data - is one long string that turns into the TUI
	var data strings.Builder

	// Clearing the terminal
	fmt.Fprintf(&data, "%s%s%s", ascii.ClearView, ascii.ClearHistory, ascii.MoveToStart)

	upperBorder := e.tui.YScroll
	lowerBorder := e.tui.YScroll + e.tui.H - 1

	// Working only with visible lines
	for i := upperBorder; i < lowerBorder; i++ {
		if i < len(e.b[e.curBuffer].Lines) {
			show := e.b[e.curBuffer].Cursor.Line() == i || e.c.ShowMD

			// This 2 variables are used to get the horizontal borders of the visible content
			start := e.tui.XScroll
			end := e.tui.W - enums.InitialOffset - len(emtpyLineSpases)

			str := e.b[e.curBuffer].Lines[i].Data
			if len(str) <= end {
				end = len(str)
			}
			if len(str) < start {
				start = 0
				end = 0
				str = []rune{}
			}

			n := tui.BuildNumber(e.b[e.curBuffer].Cursor.Line(), i+1, maxNumLen, e.c.RLN)
			var l = ""
			if e.b[e.curBuffer].IsMdFile && e.c.Render {
				switch e.curMode {
				case mode.Visual, mode.Visual_line:

					// This if statement lets me render both selected lines with highlights, and not selected with markdown render
					if (i >= e.b[e.curBuffer].Visual.Line() && i <= e.b[e.curBuffer].Cursor.Line()) || (i <= e.b[e.curBuffer].Visual.Line() && i >= e.b[e.curBuffer].Cursor.Line()) {
						l = e.b[e.curBuffer].AddVisual(e.curMode, str[start:end], i)
					} else {
						l = e.tui.BuildLine(str, show, start, end, i)
					}
				// Some other modes can use different logic for rendering, but now I just call the default for non-visual or visual_line modes
				default:
					l = e.tui.BuildLine(str, show, start, end, i)
				}

				l += ascii.Reset.Str()
			} else {
				if e.curMode == mode.Visual || e.curMode == mode.Visual_line {
					l = e.b[e.curBuffer].AddVisual(e.curMode, str[start:end], i)
				} else {
					l = string(str[start:end])
				}
			}

			// Here is where I add the line to the main data string
			fmt.Fprintf(&data, "%s %s\n\r", n, l)
		} else {
			// If the line is empty, I just add the '~' symbol
			if e.tui.ShowHello {
				fmt.Fprintf(&data, "%s%s%s\n\r", tui.Colorise("~", ascii.CyanFg), ascii.Reset, e.tui.Center(e.tui.GetASCIIInfo(i)))
			} else {
				fmt.Fprintf(&data, "%s\n\r", tui.Colorise("~", ascii.CyanFg))
			}
		}
	}

	// Calculation the visual position of the cursor
	x := e.tui.CurOff + enums.InitialOffset + len(emtpyLineSpases)
	y := e.tui.CurRow + enums.CursorLineOffset

	fmt.Fprintf(&data, "%s", ascii.Reset)

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
		fmt.Fprintf(&data, "%s", e.tui.BuildCommandBar(string(e.command)))
		fmt.Fprintf(&data, ascii.CursorBloc)

	case mode.Normal:
		var tabs []string
		for _, t := range e.b {
			tabs = append(tabs, t.Title)
		}
		cursorPos := e.tui.BuildTabs(tabs, e.curBuffer, e.c.TabNames)
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(x, y, cursorPos, e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)

	case mode.Visual, mode.Visual_line:
		fmt.Fprintf(&data, "%s", e.tui.BuildLowerBar(x, y, fmt.Sprintf("-- %s --", e.curMode), e.tui.Message, e.subCmd))
		fmt.Fprintf(&data, ascii.CursorBloc)
		fmt.Fprintf(&data, "\033[%d;%dH", y, x)
	}

	fmt.Fprintf(&data, "%s", ascii.Reset)

	// And at the end - print the data
	fmt.Print(data.String())
}
