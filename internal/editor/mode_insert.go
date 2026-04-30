package editor

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/mode"
)

// About caseInsert()
// Called from [Run()] func
// Used to insert (or delete) key in buffer
// Some specific keys (like paired symbols) can be treated differently
func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b[e.curBuffer].InsertLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = mode.Normal
		e.b[e.curBuffer].EscapeToNormal()
		e.ScrollLeft()
	case '\x7f':
		e.b[e.curBuffer].RemoveKey(0)
		e.ScrollLeft()
		e.ScrollUp()
	case '\t':
		//NOTE: yeah, I just insert 4 spaces instead of tabs
		for range 4 {
			e.b[e.curBuffer].InsertKey(' ')
			e.ScrollRight()
		}
	case '[', '{', '(', ')', '}', ']', '\'', '"', '<', '>', '*', '_', '`':
		e.b[e.curBuffer].InsertPair(key)
		e.ScrollRight()
	case 8:
		e.b[e.curBuffer].DeleteUntilSpace()
	default:
		if !isLetterNumberOrSymbol(key) {
			e.tui.Message = fmt.Sprintf("Unknown key. Code: %d", int(key))
			return
		}
		e.b[e.curBuffer].InsertKey(key)
		e.ScrollRight()
	}
	e.setUiCursor()
}

func isLetterNumberOrSymbol(key rune) bool {
	return ('a' <= key && key <= 'z') || ('A' <= key && key <= 'Z') || ('0' <= key && key <= '9') || (strings.Contains(" !@#$%^&:;|\\/~.,+=-", string(key)))
}
