package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	cases "github.com/EnotInc/Bard/internal/enums/cases"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

const visual_line = true

// Called from Run() func when current mode is Visual LIne
func (e *Editor) caseVisualLine(key rune) {
	switch key {
	case 'u':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.b[e.curBuffer].ChangeLetterCaseTo(cases.Lower, visual_line)
		e.curMode = mode.Normal

	case 'U':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.b[e.curBuffer].ChangeLetterCaseTo(cases.Upper, visual_line)
		e.curMode = mode.Normal

	case 'y':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].CopySelected(false, visual_line)
		e.curMode = mode.Normal

	case 'x':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.b[e.curBuffer].CopySelected(true, visual_line)
		e.curMode = mode.Normal

	case 'd', 'D':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.b[e.curBuffer].CopySelected(true, visual_line)
		e.curMode = mode.Normal

	case 's':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual_line)
		e.b[e.curBuffer].InsertEmptyLine(enums.Above)
		e.b[e.curBuffer].MoveToFirstChar()
		e.curMode = mode.Insert

	case 'o', 'O':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].SwapTail()

	case '<':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.execWithSubCommand(e.b[e.curBuffer].ShiftLineLeft)
		e.b[e.curBuffer].MoveToFirstVisible()

	case '>':
		if e.b[e.curBuffer].IsReadOnly {
			return
		}
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.execWithSubCommand(e.b[e.curBuffer].ShiftLineRight)
		e.b[e.curBuffer].MoveToFirstVisible()

	case keys.Esc:
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	}
}
