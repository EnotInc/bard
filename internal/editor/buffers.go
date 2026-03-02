package editor

import (
	"Enot/Bard/internal/buffer"
	"fmt"
	"slices"
)

func (e *Editor) newBuffer() {
	b := buffer.InitBuffer()
	e.b = append(e.b, b...)
	e.curBuffer = len(e.b) - 1
}

func (e *Editor) delBuffer(index int) {
	if len(e.b) > 1 {
		e.SaveFile()
		e.b = slices.Delete(e.b, index, index+1)
		e.curBuffer = 0
		e.tui.Message = fmt.Sprintf("Buffer '%s' closed", e.b[e.curBuffer].Title)
	} else {
		e.tui.Message = "Last buffer can't be removed"
	}
}
