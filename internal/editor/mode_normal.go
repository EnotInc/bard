package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/keys"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

// Called from Run() func
// used to move cursor, change move or do other stuff, depending on given key
func (e *Editor) caseNormal(key rune) {
	if ok := e.findSome(key); ok {
		return
	}
	if ok := e.replaceWith(key); ok {
		return
	}

	switch key {
	case 'r':
		e.subCmd += string(key)

	case 'i':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)
		}
		e.ScrollLeft()
		e.tui.ShowHello = false

	case 'a':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)
		}
		e.b[e.curBuffer].Insert_a()
		e.ScrollRight()
		e.tui.ShowHello = false

	case 'I':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)
		}
		e.b[e.curBuffer].MoveToFirstVisible()
		e.tui.ShowHello = false

	case 'A':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)
		}
		e.b[e.curBuffer].MoveToLastChar()
		e.b[e.curBuffer].Insert_a()
		e.tui.ShowHello = false

	case ':':
		e.curMode = mode.Command

	case 'o':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}

		e.b[e.curBuffer].InsertEmptyLine(enums.Below)
		e.b[e.curBuffer].J(1)
		e.b[e.curBuffer].Insert_a()
		e.moveToZero()

		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)
		e.tui.ShowHello = false

	case 'O':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}
		e.b[e.curBuffer].InsertEmptyLine(enums.Above)
		e.moveToZero()

		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)
		e.tui.ShowHello = false

	case 'D':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		e.b[e.curBuffer].ClearLine()
		e.b[e.curBuffer].MoveToFirstChar()
		e.moveToZero()

	case 'd':
		e.subCmd += "d"
		if e.subCmd == "dd" {
			e.b[e.curBuffer].SaveChanges(
				buffer.Delete,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(),
				enums.Without)

			e.subCmd = ""
			e.b[e.curBuffer].RemoveLine()
			e.moveToZero()
		}

	case 'R':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Replace
		}

	case 'J':
		if len(e.b[e.curBuffer].Lines)-1 == e.b[e.curBuffer].Cursor.Line() {
			return
		} // do nothing on the last line
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		line := e.b[e.curBuffer].Cursor.Line()
		e.b[e.curBuffer].J(1)

		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.With)

		e.b[e.curBuffer].DelAndMoveLineAt(line, line+1, 0)
		e.b[e.curBuffer].K(1)

	case 'x':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		e.b[e.curBuffer].Delkey()
		if e.b[e.curBuffer].Cursor.Offset() == len(e.b[e.curBuffer].Lines[e.b[e.curBuffer].Cursor.Line()].Data) {
			e.b[e.curBuffer].H(1)
		}
		e.ScrollLeft()

	case 's':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		e.b[e.curBuffer].Delkey()
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}

	case 'S':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(),
			enums.Without)

		e.b[e.curBuffer].ClearLine()
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = mode.Insert
		}

	case 'v':
		e.curMode = mode.Visual
		e.b[e.curBuffer].StartVisual()
		e.tui.ShowHello = false

	case 'V':
		e.curMode = mode.Visual_line
		e.b[e.curBuffer].StartVisualLine()
		e.tui.ShowHello = false

	case 'p':
		e.b[e.curBuffer].SaveCopied()

		e.b[e.curBuffer].Paste(enums.After)
		e.tui.ShowHello = false

	case 'P':
		e.b[e.curBuffer].SaveCopied()

		e.b[e.curBuffer].Paste(enums.Before)
		e.tui.ShowHello = false

	case 'u', keys.Ctrl_z:
		err := e.b[e.curBuffer].Undo()
		if err != nil {
			e.tui.Message = err.Error()
		}

	case keys.Ctrl_r:
		err := e.b[e.curBuffer].Redo()
		if err != nil {
			e.tui.Message = err.Error()
		}

	default:
		e.subCmd = ""
	}
}
