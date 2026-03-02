package buffer

type cursor struct {
	line   int
	offset int

	//keepOffset is uset to keep the Cursor in one place
	// on the X-axis when moving betweeen lines
	keepOffset int
}

// ============= So next 3 functions are like getters =============

func (c *cursor) Offset() int {
	return c.offset
}

func (c *cursor) Line() int {
	return c.line
}

func (c *cursor) KeepOffset() int {
	return c.keepOffset
}

// ================================================================

func (b *Buffer) FixOffset() {
	if b.Cursor.offset < b.Cursor.keepOffset {
		b.Cursor.offset = b.Cursor.keepOffset
	}
	if b.Cursor.offset > len(b.Lines[b.Cursor.line].Data)-1 {
		b.Cursor.offset = len(b.Lines[b.Cursor.line].Data) - 1
	}
	if b.Cursor.offset < 0 {
		b.Cursor.offset = 0
	}
}

// Move Cursor left
func (b *Buffer) H(amount int) {
	if b.Cursor.offset-amount > 0 {
		b.Cursor.offset -= amount
	} else {
		b.Cursor.offset = 0
	}
	b.Cursor.keepOffset = b.Cursor.offset
}

// Move Cursor down
func (b *Buffer) J(amount int) {
	if b.Cursor.line+amount < len(b.Lines)-1 {
		b.Cursor.line += amount
	} else {
		b.Cursor.line = len(b.Lines) - 1
	}
	b.FixOffset()
}

// Move Cursor up
func (b *Buffer) K(amount int) {
	if b.Cursor.line-amount > 0 {
		b.Cursor.line -= amount
	} else {
		b.Cursor.line = 0
	}
	b.FixOffset()
}

// Move Cursor right
func (b *Buffer) L(amount int) {
	if b.Cursor.offset+amount < len(b.Lines[b.Cursor.line].Data) {
		b.Cursor.offset += amount
	} else {
		b.Cursor.offset = len(b.Lines[b.Cursor.line].Data) - 1
		b.FixOffset()
	}
	b.Cursor.keepOffset = b.Cursor.offset
}

func (b *Buffer) Insert_a() {
	b.Cursor.offset += 1
}

func (b *Buffer) MoveToFirstLine() {
	b.Cursor.line = 0
	b.FixOffset()
}

func (b *Buffer) MoveToFirstChar() {
	b.Cursor.line = 0
	b.FixOffset()
}

func (b *Buffer) MoveToLastChar() {
	b.Cursor.offset = len(b.Lines[b.Cursor.line].Data)
}

func (b *Buffer) MoveToLastLine() {
	b.Cursor.line = len(b.Lines) - 1
	b.FixOffset()
}

// Move buffer.Cursor to the first non-space character in the line
func (b *Buffer) MoveToFirstVisible() {
	for i := range len(b.Lines[b.Cursor.line].Data) {
		if b.Lines[b.Cursor.line].Data[i] != ' ' {
			b.Cursor.offset = i
			break
		}
	}
	b.Cursor.keepOffset = b.Cursor.offset
}

func (b *Buffer) MoveBack(amount int) {}

func (b *Buffer) MoveBACK(amount int) {}

func (b *Buffer) MoveWord(amount int) {
	curLine := b.Lines[b.Cursor.line]
	offset := b.Cursor.offset

	for range amount {

		if len(curLine.Data) == 0 || offset == len(curLine.Data)-1 { // moving to the next line
			b.J(1)
			b.Cursor.offset = 0
			b.Cursor.keepOffset = 0
			return
		}

		ch := curLine.Data[offset]
		isSymbol := !isLetterOrNumber(ch)

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

	b.Cursor.offset = offset
	b.Cursor.keepOffset = offset
	b.FixOffset()
}

func (b *Buffer) MoveWORD(amount int) {
	curLine := b.Lines[b.Cursor.line]
	offset := b.Cursor.offset

	for range amount {

		if len(curLine.Data) == 0 || offset == len(curLine.Data)-1 { // moving to the next line
			b.J(1)
			b.Cursor.offset = 0
			b.Cursor.keepOffset = 0
			return
		}

		ch := curLine.Data[offset]

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

	b.Cursor.offset = offset
	b.Cursor.keepOffset = offset
	b.FixOffset()
}

func (b *Buffer) MoveEnd(amount int) {}

func (b *Buffer) MoveEND(amount int) {
	curLine := b.Lines[b.Cursor.line]
	offset := b.Cursor.offset

	for range amount {

		if len(curLine.Data) == 0 || offset == len(curLine.Data)-1 { // moving to the next line
			b.J(1)
			b.Cursor.offset = 0
			b.Cursor.keepOffset = 0
			return
		}

		ch := curLine.Data[offset]
		// skipping all white spaces
		if len(curLine.Data) != 0 {
			offset += 1
			ch = curLine.Data[offset]
			for offset < len(curLine.Data)-1 && ch == ' ' {
				offset += 1
				ch = curLine.Data[offset]
			}
		}

		// skipping everything until we find scpace
		for offset < len(curLine.Data)-1 && ch != ' ' {
			offset += 1
			ch = curLine.Data[offset]
		}
		for ch == ' ' { // skipping all spaces after the WORD
			offset -= 1
			ch = curLine.Data[offset]
		}
	}

	b.Cursor.offset = offset
	b.Cursor.keepOffset = offset
	b.FixOffset()
}

func isLetterOrNumber(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_' || ('0' <= ch && ch <= '9') || ch == '-'
}
