package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	cases "github.com/EnotInc/Bard/internal/enums/cases"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
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
		e.curMode = mode.Normal

	case 'x':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual)
		e.curMode = mode.Normal

	case 'u':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.b[e.curBuffer].ChangeLetterCaseTo(cases.Lower, visual)
		e.curMode = mode.Normal

	case 'U':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Visual.Line(),
			enums.Without)

		e.b[e.curBuffer].ChangeLetterCaseTo(cases.Upper, visual)
		e.curMode = mode.Normal

	case 'o', 'O':
		e.b[e.curBuffer].SwapTail()

	case 'd', 'D':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual)
		e.curMode = mode.Normal

	case 's':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, visual)
		e.curMode = mode.Insert

	case keys.Esc:
		e.curMode = mode.Normal
		e.ScrollLeft()
	}
}

func (e *Editor) saveSelected() {
	from := min(e.b[e.curBuffer].Cursor.Line(), e.b[e.curBuffer].Visual.Line())
	to := max(e.b[e.curBuffer].Cursor.Line(), e.b[e.curBuffer].Visual.Line())

	e.b[e.curBuffer].SaveChanges(
		buffer.Change,
		from, to,
		enums.Without)
	if from != to {
		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			from+1, to,
			enums.With)
	}
}
