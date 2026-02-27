package editor

import "Enot/Bard/internal/mode"

// map of paired runes. Markdown symbols are included
var openPairs map[rune]rune = map[rune]rune{
	'(':  ')',
	'[':  ']',
	'{':  '}',
	'<':  '>',
	'\'': '\'',
	'"':  '"',
	'`':  '`',
	'*':  '*',
	'_':  '_',
}

func (e *Editor) caseInsert(key rune) {
	switch key {
	case '\013', '\r', '\n':
		e.b.InsertLine()
		e.ScrollDown()
		e.moveLeft()
	case '\033':
		e.curMode = mode.Normal
		if e.b.Cursor.Offset > 0 {
			e.b.Cursor.Offset -= 1
		}
		e.pairs = []rune{}
		e.ScrollLeft()
	case '\x7f':
		e.b.RemoveKey(0)
		e.ScrollLeft()
		e.ScrollUp()
	case '\t':
		//NOTE: yeah, I just insert 4 spaces instead of tabs
		for range 4 {
			e.b.InsertKey(' ')
			e.ScrollRight()
		}
	case '[', '{', '(', ')', '}', ']', '\'', '"', '<', '>', '*', '_', '`':
		if len(e.pairs) == 0 {
			e.insertPair(key)
		} else {
			topOpen := e.pairs[len(e.pairs)-1]
			if openPairs[topOpen] == key { // if present paired key, skip pair
				e.pairs = e.pairs[:len(e.pairs)-1]
				e.b.Cursor.Offset += 1
				e.ScrollRight()
			} else {
				e.insertPair(key)
			}
		}
	default:
		e.b.InsertKey(key)
		e.ScrollRight()
	}
	e.b.Cursor.KeepOffset = e.b.Cursor.Offset
	e.setUiCursor()
}

func (e *Editor) insertPair(key rune) {
	e.b.InsertKey(key)
	if v, ok := openPairs[key]; ok { // if the key is in openPairs, insert the matching pair
		e.b.InsertKey(v)
		e.b.H(1)
		e.pairs = append(e.pairs, key)
	}
	e.ScrollRight()
}
