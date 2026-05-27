package buffer

import (
	"fmt"
	"slices"
)

type operation int

const (
	_ operation = iota
	Insert
	Change
	Delete
)

const capacity = 1000

type snapshot struct {
	lines []Line
	op    operation
	start int
	end   int
	prev  bool
}

func (b *Buffer) SaveChanges(op operation, start int, end int, with_prev bool) {
	if b.IsReadOnly {
		return
	}

	lines := []Line{}

	s := min(start, end)
	e := max(start, end)
	e += 1

	for i := range e - s {
		data := make([]rune, 0)
		data = append(data, b.Lines[i+s].Data...)
		lines = append(lines, Line{Data: data})
	}

	sn := snapshot{
		op:    op,
		lines: lines,
		start: s,
		end:   e,
		prev:  with_prev,
	}

	b.UndoStack = append(b.UndoStack, sn)

	b.fixHistoryCapacity()
	b.RedoStack = []snapshot{}
}

func (b *Buffer) fixHistoryCapacity() {
	len := len(b.UndoStack)
	if len >= capacity {
		b.UndoStack = b.UndoStack[len-capacity:]
	}
}

func (b *Buffer) Undo() error {
	if b.IsReadOnly {
		return fmt.Errorf("Buffer is read only")
	}

	if len(b.UndoStack) == 0 {
		return fmt.Errorf("Change history is empty")
	}

	index := len(b.UndoStack) - 1
	snapshot := b.UndoStack[index]

	b.Cursor.line = snapshot.start

	switch snapshot.op {
	case Change:
		snapLines := snapshot.lines

		for i, line := range snapLines {
			old_line_data := b.Lines[snapshot.start+i].Data
			b.Lines[snapshot.start+i] = &Line{Data: line.Data}
			snapshot.lines[i] = Line{Data: old_line_data}
		}
	case Insert:
		b.Lines = slices.Delete(b.Lines, snapshot.start, snapshot.end)
	case Delete:
		for i, l := range snapshot.lines {
			b.InsertLineWithData(i+snapshot.start, l.Data)
		}
	}

	b.UndoStack = b.UndoStack[:index]

	if b.Cursor.line > len(b.Lines)-1 {
		b.Cursor.line = len(b.Lines) - 1
	}
	b.MoveToFirstVisible()

	if snapshot.prev {
		err := b.Undo()
		if err != nil {
			return err
		}
	}
	b.moveToRedo(snapshot)

	return nil
}

func (b *Buffer) Redo() error {
	if b.IsReadOnly {
		return fmt.Errorf("Buffer is read only")
	}

	if len(b.RedoStack) == 0 {
		return fmt.Errorf("Change history is empty")
	}

	index := len(b.RedoStack) - 1
	snapshot := b.RedoStack[index]
	b.Cursor.line = snapshot.start

	switch snapshot.op {
	case Change:
		snapLines := snapshot.lines
		for i, line := range snapLines {
			old_line_data := b.Lines[snapshot.start+i].Data
			b.Lines[snapshot.start+i] = &Line{Data: line.Data}
			snapshot.lines[i] = Line{Data: old_line_data}
		}
	case Delete:
		b.Lines = slices.Delete(b.Lines, snapshot.start, snapshot.end)
	case Insert:
		for i, l := range snapshot.lines {
			b.InsertLineWithData(i+snapshot.start, l.Data)
		}
	}

	b.RedoStack = b.RedoStack[:index]

	if b.Cursor.line > len(b.Lines)-1 {
		b.Cursor.line = len(b.Lines) - 1
	}
	b.MoveToFirstVisible()

	if snapshot.prev {
		err := b.Redo()
		if err != nil {
			return err
		}
	}
	b.moveToUndo(snapshot)

	return nil
}

func (b *Buffer) moveToRedo(snapshot snapshot) {
	b.RedoStack = append(b.RedoStack, snapshot)
}

func (b *Buffer) moveToUndo(snapshot snapshot) {
	b.UndoStack = append(b.UndoStack, snapshot)
}
