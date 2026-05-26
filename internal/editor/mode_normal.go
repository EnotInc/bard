package editor

import (
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
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
			e.curMode = enums.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)
		}
		e.ScrollLeft()
		e.tui.ShowHello = false
	case 'a':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)
		}
		e.b[e.curBuffer].FixOffset()
		e.b[e.curBuffer].Insert_a()
		e.ScrollRight()
		e.tui.ShowHello = false
	case 'I':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)
		}
		e.b[e.curBuffer].MoveToFirstVisible()
		e.moveLeft()
		e.tui.ShowHello = false
	case 'A':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Insert
			e.b[e.curBuffer].SaveChanges(
				buffer.Change,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)
		}
		e.b[e.curBuffer].MoveToLastChar()
		e.b[e.curBuffer].Insert_a()
		e.tui.ShowHello = false
	case ':':
		e.curMode = enums.Command
	case 'o':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Insert
		}

		e.b[e.curBuffer].InsertEmptyLine(enums.Below)
		e.b[e.curBuffer].J(1)
		e.b[e.curBuffer].Insert_a()
		e.ScrollDown()
		e.moveLeft()

		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)
		e.tui.ShowHello = false
	case 'O':
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Insert
		}
		e.b[e.curBuffer].InsertEmptyLine(enums.Above)
		e.ScrollUp()
		e.moveLeft()

		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)
		e.tui.ShowHello = false
	case 'D':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		e.b[e.curBuffer].ClearLine()
		e.b[e.curBuffer].MoveToFirstChar()
		e.moveLeft()
	case 'd':
		e.subCmd += "d"
		if e.subCmd == "dd" {
			e.b[e.curBuffer].SaveChanges(
				buffer.Delete,
				e.b[e.curBuffer].Cursor.Line(),
				e.b[e.curBuffer].Cursor.Line(), false)

			e.subCmd = ""
			e.b[e.curBuffer].RemoveLine()
			e.moveLeft()
		}
	case 'R':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Replace
		}
	case 'J':
		if len(e.b[e.curBuffer].Lines)-1 == e.b[e.curBuffer].Cursor.Line() {
			return
		} // do nothing on the last line
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		line := e.b[e.curBuffer].Cursor.Line()
		e.b[e.curBuffer].J(1)

		e.b[e.curBuffer].SaveChanges(
			buffer.Delete,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), true)

		e.b[e.curBuffer].DelAndMoveLineAt(line, line+1, 0)
		e.b[e.curBuffer].K(1)
	case 'x':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		e.b[e.curBuffer].Delkey()
		if e.b[e.curBuffer].Cursor.Offset() == len(e.b[e.curBuffer].Lines[e.b[e.curBuffer].Cursor.Line()].Data) {
			e.b[e.curBuffer].H(1)
		}
		e.ScrollLeft()
	case 's':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		e.b[e.curBuffer].Delkey()
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Insert
		}
	case 'S':
		e.b[e.curBuffer].SaveChanges(
			buffer.Change,
			e.b[e.curBuffer].Cursor.Line(),
			e.b[e.curBuffer].Cursor.Line(), false)

		e.b[e.curBuffer].ClearLine()
		if !e.b[e.curBuffer].IsReadOnly {
			e.curMode = enums.Insert
		}
	case 'v':
		e.curMode = enums.Visual
		e.b[e.curBuffer].StartVisual()
		e.tui.ShowHello = false
	case 'V':
		e.curMode = enums.Visual_line
		e.b[e.curBuffer].StartVisualLine()
		e.tui.ShowHello = false
	case 'p':
		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			len(e.b[e.curBuffer].Copies)+e.b[e.curBuffer].Cursor.Line()-1, false)

		e.b[e.curBuffer].Paste(enums.After)
		e.tui.ShowHello = false
	case 'P':
		e.b[e.curBuffer].SaveChanges(
			buffer.Insert,
			e.b[e.curBuffer].Cursor.Line(),
			len(e.b[e.curBuffer].Copies)+e.b[e.curBuffer].Cursor.Line()-1, false)

		e.b[e.curBuffer].Paste(enums.Before)
		e.tui.ShowHello = false
	case 'u':
		err := e.b[e.curBuffer].Undo()
		if err != nil {
			e.tui.Message = err.Error()
		}
	case 18: // ctrl + r
		err := e.b[e.curBuffer].Redo()
		if err != nil {
			e.tui.Message = err.Error()
		}
	default:
		e.subCmd = ""
	}
}
