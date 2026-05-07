package editor

import (
	"fmt"
	"unicode"

	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
)

// About caseInsert()
// Called from [Run()] func
// Used to insert (or delete) key in buffer
// Some specific keys (like paired symbols) can be treated differently
func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		ok := e.b[e.curBuffer].DismissList()
		if ok {
			e.b[e.curBuffer].ClearLine()
			return
		}

		e.b[e.curBuffer].InsertLine()
		e.ScrollDown()
		e.moveLeft()

		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), true)

	case '\033':
		e.curMode = enums.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	case '\x7f':
		if e.b[e.curBuffer].Cursor.Offset() > 0 {
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)

			e.b[e.curBuffer].RemoveKey()
		} else {
			e.b[e.curBuffer].SaveChanges(
				buffer.Delete,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)
			e.b[e.curBuffer].K(1)
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), true)
			e.b[e.curBuffer].J(1)

			e.b[e.curBuffer].DelAndMoveLine()
		}

		e.ScrollLeft()
		e.ScrollUp()

	case '\t':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		//NOTE: yeah, I just insert 4 spaces instead of tabs
		for range 4 {
			e.b[e.curBuffer].InsertKey(' ')
			e.ScrollRight()
		}
	case '[', '{', '(', ')', '}', ']', '\'', '"', '<', '>', '*', '_', '`':
		e.b[e.curBuffer].InsertPair(key)
		e.ScrollRight()
	case 8:
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		e.b[e.curBuffer].DeleteUntilSpace()
	default:
		if key == ' ' {
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)
		}
		if !unicode.IsPrint(key) {
			e.tui.Message = fmt.Sprintf("Unknown key. Code: %d", int(key))
			return
		}
		e.b[e.curBuffer].InsertKey(key)
		e.ScrollRight()
	}
	e.setUiCursor()
}
