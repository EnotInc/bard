package editor

const (
	ScrollBorder     = 5
	cursorLineOffset = 1
	initialOffset    = 3
)

// Making sure that visual Cursor is alright
func (e *Editor) setUiCursor() {
	if e.tui.XScroll > e.b.Cursor.Offset() {
		e.tui.XScroll = e.b.Cursor.Offset()
	} else if e.tui.XScroll < e.b.Cursor.Offset()-e.tui.W+ScrollBorder*2 {
		e.tui.XScroll = e.b.Cursor.Offset() - e.tui.W + ScrollBorder*2
	}
	if e.tui.YScroll > e.b.Cursor.Line()-ScrollBorder {
		e.tui.YScroll = e.b.Cursor.Line() - ScrollBorder
		if e.tui.YScroll < 0 {
			e.tui.YScroll = 0
		}
	} else if e.tui.YScroll < e.b.Cursor.Line()-e.tui.H+ScrollBorder {
		e.tui.YScroll = e.b.Cursor.Line() - e.tui.H + ScrollBorder
	}

	e.tui.CurRow = e.b.Cursor.Line() - e.tui.YScroll
	e.tui.CurOff = e.b.Cursor.Offset() - e.tui.XScroll
}

func (e *Editor) ScrollUp() {
	if e.tui.CurRow == ScrollBorder {
		if e.tui.YScroll != 0 {
			e.tui.YScroll -= 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollDown() {
	if e.tui.CurRow == e.tui.H-ScrollBorder {
		if e.tui.YScroll+e.tui.H != len(e.b.Lines)+ScrollBorder {
			e.tui.YScroll += 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollRight() {
	if e.tui.CurOff >= e.tui.W-ScrollBorder*2 {
		if e.tui.XScroll+e.tui.W-initialOffset != len(e.b.Lines[e.b.Cursor.Line()].Data)+ScrollBorder {
			e.tui.XScroll += 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollLeft() {
	if e.tui.CurOff <= ScrollBorder {
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
	e.tui.XScroll = len(e.b.Lines[e.b.Cursor.Line()].Data) - e.tui.W + ScrollBorder*2
	if e.tui.XScroll < 0 {
		e.tui.XScroll = 0
	}
	e.b.ResetKeepOffset()
	e.setUiCursor()
}

func (e *Editor) shiftLeft() {
	if e.tui.XScroll > e.b.Cursor.Offset() {
		e.tui.XScroll = e.b.Cursor.Offset() - ScrollBorder
		if e.tui.XScroll < 0 {
			e.tui.XScroll = 0
		}
	}
	e.setUiCursor()
}
