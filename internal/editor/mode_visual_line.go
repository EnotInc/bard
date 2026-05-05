package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/mode"
)

// About caseVisualLine()
// Called from [Run()] func when current mode is Visual LIne
func (e *Editor) caseVisualLine(key rune) {
	switch key {
	case 'y':
		e.b[e.curBuffer].CopySelected(false, true)
		e.curMode = mode.Normal
	case 'x':
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].CopySelected(true, true)
		e.curMode = mode.Normal
	case 'd', 'D':
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].CopySelected(true, true)
		e.curMode = mode.Normal
	case 's':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, true)
		e.b[e.curBuffer].InsertEmptyLine(enums.Above)
		e.b[e.curBuffer].MoveToFirstChar()
		e.curMode = mode.Insert
	case 'o', 'O':
		e.b[e.curBuffer].SwapTail()
	case '\033':
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	}
}
