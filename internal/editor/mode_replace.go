package editor

import (
	"unicode"

	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

// Unless `escape` key is pressed, replaces 1 chat (given key) at the time
func (e *Editor) caseReplaceChar(key rune, amount int) {
	switch key {
	case keys.Enter:
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		e.b[e.curBuffer].Delkey()
		e.b[e.curBuffer].InsertLine()
		e.moveToZero()

		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.With)

	case keys.Esc:
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()

	case keys.Backspace:
		return

	case keys.Tab:
		e.b[e.curBuffer].Delkey()
		e.b[e.curBuffer].InsertKey('\t')

		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)
	default:
		if unicode.IsPrint(key) {
			e.b[e.curBuffer].ReplaceKeys(key, amount)
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)
		}
	}
	e.subCmd = ""
}

// Called from Run() func
// Untin `escape` key is pressed, replaces char and moves curosr forward
func (e *Editor) caseReplaceMode(key rune) {
	switch key {
	case keys.Enter:
		e.b[e.curBuffer].DelAndMoveLine()
		e.moveToZero()

	case keys.Esc:
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()

	case keys.Backspace:
		e.b[e.curBuffer].RemoveKey()
		e.ScrollLeft()

	case keys.Tab:
		e.b[e.curBuffer].Delkey()
		e.b[e.curBuffer].InsertKey('\t')

	default:
		if unicode.IsPrint(key) {
			e.b[e.curBuffer].ReplaceKeys(key, 1)
			e.b[e.curBuffer].L(1)
		}
	}
}
