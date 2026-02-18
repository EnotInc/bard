package editor

const (
	ScrollBorder = 5
)

// Making sure that visual cursor is alright
func (e *Editor) setUiCursor() {
	if e.ui.xScroll > e.b.cursor.offset {
		e.ui.xScroll = e.b.cursor.offset
	} else if e.ui.xScroll < e.b.cursor.offset-e.ui.w+ScrollBorder*2 {
		e.ui.xScroll = e.b.cursor.offset - e.ui.w + ScrollBorder*2
	}
	if e.ui.yScroll > e.b.cursor.line-ScrollBorder {
		e.ui.yScroll = e.b.cursor.line - ScrollBorder
		if e.ui.yScroll < 0 {
			e.ui.yScroll = 0
		}
	} else if e.ui.yScroll < e.b.cursor.line-e.ui.h+ScrollBorder {
		e.ui.yScroll = e.b.cursor.line - e.ui.h + ScrollBorder
	}

	e.ui.curRow = e.b.cursor.line - e.ui.yScroll
	e.ui.curOff = e.b.cursor.offset - e.ui.xScroll
}

func (e *Editor) ScrollUp() {
	if e.ui.curRow == ScrollBorder {
		if e.ui.yScroll != 0 {
			e.ui.yScroll -= 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollDown() {
	if e.ui.curRow == e.ui.h-ScrollBorder {
		if e.ui.yScroll+e.ui.h != len(e.b.lines)+ScrollBorder {
			e.ui.yScroll += 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollRight() {
	if e.ui.curOff >= e.ui.w-ScrollBorder*2 {
		if e.ui.xScroll+e.ui.w-initialOffset != len(e.b.lines[e.b.cursor.line].data)+ScrollBorder {
			e.ui.xScroll += 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) ScrollLeft() {
	if e.ui.curOff <= ScrollBorder {
		if e.ui.xScroll != 0 {
			e.ui.xScroll -= 1
		}
	}
	e.setUiCursor()
}

func (e *Editor) moveLeft() {
	e.ui.xScroll = 0
	e.setUiCursor()
}

func (e *Editor) moveRight() {
	e.ui.xScroll = len(e.b.lines[e.b.cursor.line].data) - e.ui.w + ScrollBorder*2
	if e.ui.xScroll < 0 {
		e.ui.xScroll = 0
	}
	e.setUiCursor()
}

func (e *Editor) shiftLeft() {
	if e.ui.xScroll > e.b.cursor.offset {
		e.ui.xScroll = e.b.cursor.offset - ScrollBorder
		if e.ui.xScroll < 0 {
			e.ui.xScroll = 0
		}
	}
	e.setUiCursor()
}
