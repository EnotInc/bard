package editor

import (
	"unicode"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/services"
)

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
		buf := e.b[e.curBuffer]
		buf.Delkey()

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

		buf.SaveChanges(
			buffer.Change,
			buf.Cursor.Line(),
			buf.Cursor.Line(),
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
		buf := e.b[e.curBuffer]
		buf.Delkey()

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

	default:
		if unicode.IsPrint(key) {
			e.b[e.curBuffer].ReplaceKeys(key, 1)
			e.b[e.curBuffer].L(1)
		}
	}
}
