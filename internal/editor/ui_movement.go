package editor

import (
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/services"
)

// Making sure that visual Cursor is not out of bounds
func (e *Editor) setUiCursor() {
	buf := e.b[e.curBuffer]

	if e.tui.XScroll > buf.Cursor.Offset() {
		e.tui.XScroll = buf.Cursor.Offset()
	} else if e.tui.XScroll < buf.Cursor.Offset()-e.tui.W+enums.ScrollBorder*2 {
		e.tui.XScroll = buf.Cursor.Offset() - e.tui.W + enums.ScrollBorder*2
	}
	if e.tui.YScroll > buf.Cursor.Line()-enums.ScrollBorder {
		e.tui.YScroll = buf.Cursor.Line() - enums.ScrollBorder
		if e.tui.YScroll < 0 {
			e.tui.YScroll = 0
		}
	} else if e.tui.YScroll < buf.Cursor.Line()-e.tui.H+enums.ScrollBorder {
		e.tui.YScroll = buf.Cursor.Line() - e.tui.H + enums.ScrollBorder
	}

	shift := services.CursorShiftAt(buf.Lines[buf.Cursor.Line()].Data, buf.Cursor.Offset())
	e.tui.CurOff = buf.Cursor.Offset() - e.tui.XScroll + shift
	e.tui.CurRow = buf.Cursor.Line() - e.tui.YScroll
}

// changes YScroll in TUI by decreaseing it, only if CurRow is equal to ScrollBorder
func (e *Editor) ScrollUp() {
	if e.tui.CurRow == enums.ScrollBorder {
		if e.tui.YScroll != 0 {
			e.tui.YScroll -= 1
		}
	}
}

// changes YScroll in TUI by adding one to it, only if CurRow is equal to ScrollBorder (at the bottom)
func (e *Editor) ScrollDown() {
	if e.tui.CurRow == e.tui.H-enums.ScrollBorder {
		if e.tui.YScroll+e.tui.H != len(e.b[e.curBuffer].Lines)+enums.ScrollBorder {
			e.tui.YScroll += 1
		}
	}
}

// changes XScroll if CurOff is toching ScrollBorder
func (e *Editor) ScrollRight() {
	if e.tui.CurOff >= e.tui.W-enums.ScrollBorder*2 {
		if e.tui.XScroll+e.tui.W-enums.InitialOffset != len(e.b[e.curBuffer].Lines[e.b[e.curBuffer].Cursor.Line()].Data)+enums.ScrollBorder {
			e.tui.XScroll += 1
		}
	}
}

// changes XScroll if CurOff is toching ScrollBorder
func (e *Editor) ScrollLeft() {
	if e.tui.CurOff <= enums.ScrollBorder {
		if e.tui.XScroll != 0 {
			e.tui.XScroll -= 1
		}
	}
}

// About moveLeft()
// sets XScroll to zero
func (e *Editor) moveLeft() {
	e.tui.XScroll = 0
}

// shift XScroll to the left
// moves to the first non space char in line
func (e *Editor) shiftLeft() {
	if e.tui.XScroll > e.b[e.curBuffer].Cursor.Offset() {
		e.tui.XScroll = max(e.b[e.curBuffer].Cursor.Offset()-enums.ScrollBorder, 0)
	}
}
