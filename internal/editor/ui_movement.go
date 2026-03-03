package editor

import "Enot/Bard/internal/enums"

// Making sure that visual Cursor is alright
func (e *Editor) setUiCursor() {
	if e.tui.XScroll > e.b[e.curBuffer].Cursor.Offset() {
		e.tui.XScroll = e.b[e.curBuffer].Cursor.Offset()
	} else if e.tui.XScroll < e.b[e.curBuffer].Cursor.Offset()-e.tui.W+enums.ScrollBorder*2 {
		e.tui.XScroll = e.b[e.curBuffer].Cursor.Offset() - e.tui.W + enums.ScrollBorder*2
	}
	if e.tui.YScroll > e.b[e.curBuffer].Cursor.Line()-enums.ScrollBorder {
		e.tui.YScroll = e.b[e.curBuffer].Cursor.Line() - enums.ScrollBorder
		if e.tui.YScroll < 0 {
			e.tui.YScroll = 0
		}
	} else if e.tui.YScroll < e.b[e.curBuffer].Cursor.Line()-e.tui.H+enums.ScrollBorder {
		e.tui.YScroll = e.b[e.curBuffer].Cursor.Line() - e.tui.H + enums.ScrollBorder
	}

	e.tui.CurRow = e.b[e.curBuffer].Cursor.Line() - e.tui.YScroll
	e.tui.CurOff = e.b[e.curBuffer].Cursor.Offset() - e.tui.XScroll
}

func (e *Editor) ScrollUp() {
	if e.tui.CurRow == enums.ScrollBorder {
		if e.tui.YScroll != 0 {
			e.tui.YScroll -= 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollDown() {
	if e.tui.CurRow == e.tui.H-enums.ScrollBorder {
		if e.tui.YScroll+e.tui.H != len(e.b[e.curBuffer].Lines)+enums.ScrollBorder {
			e.tui.YScroll += 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollRight() {
	if e.tui.CurOff >= e.tui.W-enums.ScrollBorder*2 {
		if e.tui.XScroll+e.tui.W-enums.InitialOffset != len(e.b[e.curBuffer].Lines[e.b[e.curBuffer].Cursor.Line()].Data)+enums.ScrollBorder {
			e.tui.XScroll += 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollLeft() {
	if e.tui.CurOff <= enums.ScrollBorder {
		if e.tui.XScroll != 0 {
			e.tui.XScroll -= 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) moveLeft() {
	e.tui.XScroll = 0
	e.setUiCursor()
}

func (e *Editor) moveRight() {
	e.tui.XScroll = len(e.b[e.curBuffer].Lines[e.b[e.curBuffer].Cursor.Line()].Data) - e.tui.W + enums.ScrollBorder*2
	if e.tui.XScroll < 0 {
		e.tui.XScroll = 0
	}
	e.b[e.curBuffer].ResetKeepOffset()
	e.setUiCursor()
}

func (e *Editor) shiftLeft() {
	if e.tui.XScroll > e.b[e.curBuffer].Cursor.Offset() {
		e.tui.XScroll = e.b[e.curBuffer].Cursor.Offset() - enums.ScrollBorder
		if e.tui.XScroll < 0 {
			e.tui.XScroll = 0
		}
	}
	e.setUiCursor()
}
