package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
)

const visual_line = true

// About caseVisualLine()
// Called from [Run()] func when current mode is Visual LIne
func (e *Editor) caseVisualLine(key rune) {
	switch key {
	case 'u':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].ChangeLetterCaseTo(enums.Lower, visual_line)
		e.curMode = enums.Normal
	case 'U':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].ChangeLetterCaseTo(enums.Upper, visual_line)
		e.curMode = enums.Normal
	case 'y':
		e.b[e.curBuffer].CopySelected(false, visual_line)
		e.curMode = enums.Normal
	case 'x':
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].CopySelected(true, visual_line)
		e.curMode = enums.Normal
	case 'd', 'D':
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].CopySelected(true, visual_line)
		e.curMode = enums.Normal
	case 's':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual_line)
		e.b[e.curBuffer].InsertEmptyLine(enums.Above)
		e.b[e.curBuffer].MoveToFirstChar()
		e.curMode = enums.Insert
	case 'o', 'O':
		e.b[e.curBuffer].SwapTail()
	case '<':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.execWithSubCommand(e.b[e.curBuffer].ShiftLineLeft)
		e.b[e.curBuffer].MoveToFirstVisible()
	case '>':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.execWithSubCommand(e.b[e.curBuffer].ShiftLineRight)
		e.b[e.curBuffer].MoveToFirstVisible()
	case '\033':
		e.curMode = enums.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	}
}
