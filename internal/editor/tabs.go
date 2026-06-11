package editor

import (
	"fmt"
	"slices"

	"github.com/EnotInc/Bard/internal/editor/buffer"
)

func (e *Editor) ClearAllBuffers() {
	e.curBuffer = 0
	e.b = buffer.InitBuffer()
}

// used to create new buffer, and switch to it
func (e *Editor) newBuffer() {
	b := buffer.InitBuffer()
	e.b = append(e.b, b...)
	e.curBuffer = len(e.b) - 1
}

// deletes buffer by given index, unless current buffer is the last one
func (e *Editor) delBuffer(index int) {
	if len(e.b) > 1 {
		title := e.b[index].Title
		e.b = slices.Delete(e.b, index, index+1)
		e.curBuffer = 0
		e.tui.Message = fmt.Sprintf("Buffer '%s' closed", title)
	} else {
		e.tui.Message = "Last buffer can't be removed"
	}
}

// changes tabs, by increasing curBuffer
func (e *Editor) nextTab() {
	if e.curBuffer+1 >= len(e.b) {
		e.curBuffer = 0
	} else {
		e.curBuffer += 1
	}
}

// changes tabs, by decreasing curBuffer
// similar to nextTab()
func (e *Editor) prevTab() {
	if e.curBuffer-1 <= 0 {
		e.curBuffer = len(e.b) - 1
	} else {
		e.curBuffer -= 1
	}
}

func (e *Editor) SwitchToTab(index int) {
	if index > len(e.b)-1 || index < 0 {
		e.tui.Message = fmt.Sprintf("Unable to open tab, index %d out of range", index)
		return
	}

	e.curBuffer = index
}
