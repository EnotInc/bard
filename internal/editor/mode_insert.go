package editor

import (
	"fmt"
	"unicode"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/editor/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/services"
)

// Called from Run() func
// Used to insert (or delete) key in buffer
// Some specific keys (like paired symbols) can be treated differently
func (e *Editor) caseInsert(key rune) {
	switch key {
	case keys.Enter:
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		ok := e.b[e.curBuffer].DismissList()
		if ok {
			e.b[e.curBuffer].ClearLine()
			return
		}

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
		if e.b[e.curBuffer].Cursor.Offset() > 0 {
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)

			e.b[e.curBuffer].RemoveKey()
		} else {
			e.b[e.curBuffer].SaveChanges(
				buffer.Delete,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)

			e.b[e.curBuffer].K(1)
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.With)

			e.b[e.curBuffer].J(1)
			e.b[e.curBuffer].DelAndMoveLine()
		}

		e.ScrollLeft()

	case keys.Tab:
		buf := e.b[e.curBuffer]
		buf.SaveChanges(
			buffer.Change,
			buf.Cursor.Line(),
			buf.Cursor.Line(),
			enums.Without)

		cfg := config.GetConfig()
		if !cfg.KeepTabs {
			curLine := buf.Lines[buf.Cursor.Line()]
			tab := services.CursorShiftCalculateAt(curLine.Data, buf.Cursor.Offset())
			for range tab {
				buf.InsertKey(' ')
			}
		} else {
			buf.InsertKey('\t')
		}

	case '[', '{', '(', ')', '}', ']', '\'', '"', '<', '>', '*', '_', '`':
		e.b[e.curBuffer].InsertPair(key)
		e.ScrollRight()

	case keys.Ctrl_Backspace:
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		e.b[e.curBuffer].DeleteUntilSpace()

	default:
		if key == keys.Space {
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)
		}
		if !unicode.IsPrint(key) {
			e.tui.Error = fmt.Sprintf("Unknown key. Code: %d", int(key))
			return
		}
		e.b[e.curBuffer].InsertKey(key)
		e.ScrollRight()
	}
}
