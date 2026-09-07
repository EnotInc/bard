package editor

import (
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/services"
)

// FIXME: move to screen?
func (e *Editor) fromKeyMap(r rune) rune {
	if services.IsLetter(r) || e.curMode == mode.Insert {
		return r
	}

	if len(e.keyMaps) == 0 {
		return r
	}

	for _, m := range e.keyMaps {
		if key, ok := m[r]; ok {
			return key
		}
	}

	return r
}
