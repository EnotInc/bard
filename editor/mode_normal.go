package editor

const (
	above = iota
	below
)

func (e *Editor) caseNormal(key rune) {
	switch key {
	case 'h':
		e.b.H()
		e.ScrollLeft()
	case 'j':
		e.b.J()
		e.ScrollDown()
		e.shiftLeft()
	case 'k':
		e.b.K()
		e.ScrollUp()
		e.shiftLeft()
	case 'l':
		e.b.L()
		e.ScrollRight()
	case 'v':
		e.curMode = visual
	case 'i':
		e.curMode = insert
		e.ScrollLeft()
	case 'a':
		e.curMode = insert
		if len(e.b.lines[e.b.cursor.line].data) > 0 {
			e.b.cursor.ofset += 1
		}
		e.ScrollRight()
	case 'I':
		e.curMode = insert
		e.b.moveToFirst()
		e.moveLeft()
	case 'A':
		e.curMode = insert
		e.b.cursor.ofset = len(e.b.lines[e.b.cursor.line].data)
		e.moveRight()
	case ':':
		e.curMode = command
	case 'o':
		e.curMode = insert
		e.b.cursor.ofset = 0
		e.b.InsertEmptyLine(below)
		e.b.cursor.line += 1
		e.ScrollDown()
		e.moveLeft()
	case 'O':
		e.curMode = insert
		e.b.cursor.ofset = 0
		e.b.InsertEmptyLine(above)
		e.ScrollUp()
		e.moveLeft()
	case 'x':
		//e.b.Yank()
		e.b.Delkey()
		if e.b.cursor.ofset >= len(e.b.lines[e.b.cursor.line].data) && e.b.cursor.ofset > 0 {
			e.b.cursor.ofset -= 1
		}
		e.ScrollLeft()
	case 's':
		e.b.Delkey()
		e.curMode = insert
	}
	e.setUiCursor()
}
