package buffer

import (
	"slices"
)

type Buffer struct {
	copies []*copied
	Lines  []*Line
	Cursor *cursor
	Visual *cursor
}

func InitBuffer() *Buffer {
	c := &cursor{Line: 0, Offset: 0}
	v := &cursor{Line: 0, Offset: 0}
	b := &Buffer{
		Cursor: c,
		Visual: v,
	}
	b.Lines = append(b.Lines, &Line{Data: []rune("")})
	return b
}

func (b *Buffer) InsertKey(key rune) {
	curLine := b.Lines[b.Cursor.Line]
	curLine.Data = append(curLine.Data[:b.Cursor.Offset], append([]rune{key}, curLine.Data[b.Cursor.Offset:]...)...)
	b.Cursor.Offset += 1
}

func (b *Buffer) ReplaceKeys(key rune, amount int) {
	curLine := b.Lines[b.Cursor.Line]
	if b.Cursor.Offset < len(curLine.Data) {
		if b.Cursor.Offset+amount <= len(curLine.Data) {
			curLine.Data = slices.Delete(curLine.Data, b.Cursor.Offset, b.Cursor.Offset+amount-1)
			curLine.Data[b.Cursor.Offset] = key
		}
	} else {
		b.InsertKey(key)
	}
	b.fixOffset()
}

// Called when the user presses [backspace] and just removes the character in front of it
func (b *Buffer) RemoveKey(keyShift int) {
	if b.Cursor.Offset > 0 {
		curLine := b.Lines[b.Cursor.Line]
		index := keyShift + b.Cursor.Offset
		curLine.Data = slices.Delete(curLine.Data, index-1, index)
		b.Cursor.Offset -= 1
	} else {
		b.DelAndMoveLine()
	}
}

// Called when the user presses [x] or [s] in normal mode. It deletes the character and copies it to the buffer
func (b *Buffer) Delkey() {
	if len(b.Lines[b.Cursor.Line].Data) > 0 {
		curLine := b.Lines[b.Cursor.Line]
		index := b.Cursor.Offset
		ch := curLine.Data[index]
		b.copies = append([]*copied{}, &copied{data: []rune{ch}, isStart: false, isEnd: false})
		curLine.Data = slices.Delete(curLine.Data, index, index+1)
	}
}

func (b *Buffer) MoveBack(amount int) {}

func (b *Buffer) MoveWord(amount int) {
	curLine := b.Lines[b.Cursor.Line]
	offset := b.Cursor.Offset
	ch := curLine.Data[offset]
	isSymbol := !isLetterOrNumber(ch)

	for range amount {
		if offset == len(curLine.Data)-1 { // moving to the next line
			b.J(1)
			b.Cursor.Offset = 0
			b.Cursor.KeepOffset = 0
			return
		}

		if isSymbol {
			symbol := ch
			for symbol == ch && offset < len(curLine.Data)-1 {
				offset += 1
				ch = curLine.Data[offset]
			}
		} else {
			for offset < len(curLine.Data)-1 && isLetterOrNumber(ch) && ch != ' ' {
				offset += 1
				ch = curLine.Data[offset]
			}
		}
		for ch == ' ' { // skipping all spaces after the word
			offset += 1
			ch = curLine.Data[offset]
		}
	}

	b.Cursor.Offset = offset
	b.Cursor.KeepOffset = offset
	b.fixOffset()
}

func (b *Buffer) MoveWORD(amount int) {
	curLine := b.Lines[b.Cursor.Line]
	offset := b.Cursor.Offset
	ch := curLine.Data[offset]

	for range amount {
		if offset == len(curLine.Data)-1 { // moving to the next line
			b.J(1)
			b.Cursor.Offset = 0
			b.Cursor.KeepOffset = 0
			return
		}

		// skipping everything until we find scpace
		for offset < len(curLine.Data)-1 && ch != ' ' {
			offset += 1
			ch = curLine.Data[offset]
		}
		for ch == ' ' { // skipping all spaces after the WORD
			offset += 1
			ch = curLine.Data[offset]
		}
	}

	b.Cursor.Offset = offset
	b.Cursor.KeepOffset = offset
	b.fixOffset()
}

func (b *Buffer) MoveEnd(amount int) {
	curLine := b.Lines[b.Cursor.Line]
	offset := b.Cursor.Offset
	if len(curLine.Data) == 0 || offset == len(curLine.Data)-1 {
		b.J(1)
		b.Cursor.Offset = 0
		b.Cursor.KeepOffset = 0
		return
	}

	ch := curLine.Data[offset]

	if offset+1 < len(curLine.Data)-1 {
		offset += 1
		ch = curLine.Data[offset]
		for ch == ' ' && offset < len(curLine.Data)-1 {
			offset += 1
			ch = curLine.Data[offset]
		}
	}

	isSymbol := !isLetterOrNumber(ch) && ch != ' '

	if isSymbol {
		symbol := ch
		for (ch == symbol || ch == ' ') && offset < len(curLine.Data)-1 {
			offset += 1
			ch = curLine.Data[offset]
		}
		if offset > 1 {
			offset -= 1
		}
	} else {
		for isLetterOrNumber(ch) && offset < len(curLine.Data)-1 {
			offset += 1
			ch = curLine.Data[offset]
		}
		for ch == ' ' || !isLetterOrNumber(ch) {
			offset -= 1
			ch = curLine.Data[offset]
		}
	}

	b.Cursor.Offset = offset
	b.Cursor.KeepOffset = offset
	b.fixOffset()
}

func isLetterOrNumber(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_' || ('0' <= ch && ch <= '9') || ch == '-'
}
