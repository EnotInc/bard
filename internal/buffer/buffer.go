package buffer

import (
	"slices"
)

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

type Buffer struct {
	Title    string
	pairs    []rune // paired brackets
	copies   []*copied
	Lines    []*Line
	Cursor   *cursor
	Visual   *cursor
	IsMdFile bool
}

func InitBuffer() []*Buffer {
	c := &cursor{line: 0, offset: 0}
	v := &cursor{line: 0, offset: 0}
	b := &Buffer{
		Cursor:   c,
		Visual:   v,
		pairs:    []rune{},
		IsMdFile: false,
	}
	b.Lines = append(b.Lines, &Line{Data: []rune("")})
	var bfs []*Buffer
	bfs = append(bfs, b)
	return bfs
}

func (b *Buffer) InsertKey(key rune) {
	curLine := b.Lines[b.Cursor.line]
	curLine.Data = append(curLine.Data[:b.Cursor.offset], append([]rune{key}, curLine.Data[b.Cursor.offset:]...)...)
	b.Cursor.offset += 1
}

func (b *Buffer) ReplaceKeys(key rune, amount int) {
	curLine := b.Lines[b.Cursor.line]
	if b.Cursor.offset < len(curLine.Data) {
		if b.Cursor.offset+amount <= len(curLine.Data) {
			curLine.Data = slices.Delete(curLine.Data, b.Cursor.offset, b.Cursor.offset+amount-1)
			curLine.Data[b.Cursor.offset] = key
		}
	} else {
		b.InsertKey(key)
	}
	b.FixOffset()
}

// Called when the user presses [backspace] and just removes the character in front of it
func (b *Buffer) RemoveKey(keyShift int) {
	if b.Cursor.offset > 0 {
		curLine := b.Lines[b.Cursor.line]
		index := keyShift + b.Cursor.offset
		curLine.Data = slices.Delete(curLine.Data, index-1, index)
		b.Cursor.offset -= 1
	} else {
		b.DelAndMoveLine()
	}
}

// Called when the user presses [x] or [s] in normal mode. It deletes the character and copies it to the buffer
func (b *Buffer) Delkey() {
	if len(b.Lines[b.Cursor.line].Data) > 0 {
		curLine := b.Lines[b.Cursor.line]
		index := b.Cursor.offset
		ch := curLine.Data[index]
		b.copies = append([]*copied{}, &copied{data: []rune{ch}, isStart: false, isEnd: false})
		curLine.Data = slices.Delete(curLine.Data, index, index+1)
	}
}
func (b *Buffer) InsertPair(key rune) {
	if len(b.pairs) == 0 {
		b.insertBoth(key)
	} else {
		topOpen := b.pairs[len(b.pairs)-1]
		if openPairs[topOpen] == key { // if present paired key, skip pair
			b.pairs = b.pairs[:len(b.pairs)-1]
			b.Cursor.offset += 1
		} else {
			b.insertBoth(key)
		}
	}
}

func (b *Buffer) insertBoth(key rune) {
	if v, ok := openPairs[key]; ok { // if the key is in openPairs, insert the matching pair
		b.InsertKey(key)
		b.InsertKey(v)
		b.H(1)
		b.pairs = append(b.pairs, key)
	}
}

func (b *Buffer) EscapeToNormal() {
	b.pairs = []rune{}
	if b.Cursor.offset > 0 {
		b.Cursor.offset -= 1
	}
}

func (b *Buffer) ResetKeepOffset() {
	b.Cursor.keepOffset = b.Cursor.offset
}
