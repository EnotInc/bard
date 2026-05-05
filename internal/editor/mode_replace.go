package editor

import (
	"unicode"

	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/mode"
)

// About caseReplaceChar()
// Unless `escape` key is pressed, replaces 1 chat (given key) at the time
func (e *Editor) caseReplaceChar(key rune, amount int) {
	switch key {
	case '\013', '\r', '\n':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		e.b[e.curBuffer].InsertLine()
		e.ScrollDown()
		e.moveLeft()

		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), true)
	case '\033':
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
		e.setUiCursor()
	case '\x7f': // just do nothing if backspace is pressed
		return
	case '\t':
		e.b[e.curBuffer].Delkey()
		for range 4 {
			e.b[e.curBuffer].InsertKey(' ')
			e.ScrollRight()
		}
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)
	default:
		if unicode.IsPrint(key) {
			e.b[e.curBuffer].ReplaceKeys(key, amount)
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)
		}
	}
	e.subCmd = ""
}

// About caseReplaceMode()
// Called from [Run()] func
// Untin `escape` key is pressed, replaces char and moves curosr forward
func (e *Editor) caseReplaceMode(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b[e.curBuffer].DelAndMoveLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	case '\x7f':
		e.b[e.curBuffer].RemoveKey(0)
		e.ScrollLeft()
		e.ScrollUp()
	case '\t':
		e.b[e.curBuffer].Delkey()
		for range 4 {
			e.b[e.curBuffer].InsertKey(' ')
			e.ScrollRight()
		}
	default:
		if unicode.IsPrint(key) {
			e.b[e.curBuffer].ReplaceKeys(key, 1)
			e.b[e.curBuffer].L(1)
		}
	}
}
