package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
)

const visual = false

// kinda similar to caseVisulLine
func (e *Editor) caseVisual(key rune) {
	if ok := e.findSome(key); ok {
		return
	}
	if ok := e.replaceWith(key); ok {
		return
	}

	switch key {
	case 'y':
		e.b[e.curBuffer].CopySelected(false, visual)
		e.curMode = enums.Normal
	case 'x':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual)
		e.curMode = enums.Normal
	case 'u':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].ChangeLetterCaseTo(enums.Lower, visual)
		e.curMode = enums.Normal
	case 'U':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(), false)

		e.b[e.curBuffer].ChangeLetterCaseTo(enums.Upper, visual)
		e.curMode = enums.Normal
	case 'o', 'O':
		e.b[e.curBuffer].SwapTail()
	case 'd', 'D':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual)
		e.curMode = enums.Normal
	case 's':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual)
		e.curMode = enums.Insert
	case '\033':
		e.curMode = enums.Normal
		e.ScrollLeft()
	}
}

func (e *Editor) saveSelected() {
	from := min(e.b[e.curBuffer].Cursor.Line(), e.b[e.curBuffer].Visual.Line())
	to := max(e.b[e.curBuffer].Cursor.Line(), e.b[e.curBuffer].Visual.Line())

	e.b[e.curBuffer].SaveChanges(
		buffer.Change,
		from, to, false)
	if from != to {
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			from+1, to, true)
	}
}
