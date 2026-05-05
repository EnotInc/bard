package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/mode"
)

// About caseVisual()
// kinda similar to caseVisulLine
func (e *Editor) caseVisual(key rune) {
	cmd := []byte(e.subCmd)
	if len(cmd) > 0 {
		switch cmd[len(cmd)-1] {
		case 'f':
			e.b[e.curBuffer].FindNext(key)
			e.subCmd = ""
			return
		case 'F':
			e.b[e.curBuffer].FindPrev(key)
			e.subCmd = ""
			return
		case 't':
			e.b[e.curBuffer].FindBeforeNext(key)
			e.subCmd = ""
			return
		case 'T':
			e.b[e.curBuffer].FindBeforePrev(key)
			e.subCmd = ""
			return
		}
	}

	switch key {
	case 'y':
		e.b[e.curBuffer].CopySelected(false, false)
		e.curMode = mode.Normal
	case 'x':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, false)
		e.curMode = mode.Normal
	case 'o', 'O':
		e.b[e.curBuffer].SwapTail()
	case 'd', 'D':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, false)
		e.curMode = mode.Normal
	case 's':
		e.saveSelected()

		e.b[e.curBuffer].CopySelected(true, false)
		e.curMode = mode.Insert
	case '\033':
		e.curMode = mode.Normal
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
