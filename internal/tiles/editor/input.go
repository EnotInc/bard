package editor

import "github.com/EnotInc/Bard/internal/services/text"

func (e *Editor) fromKeyMap(r rune) rune {
	if text.IsLetter(r) {
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
