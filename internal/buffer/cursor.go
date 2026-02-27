package buffer

type cursor struct {
	Line   int
	Offset int

	//keepOffset is uset to keep the cCursor in one place
	// on the X-axis when moving betweeen lines
	KeepOffset int
}

func (b *Buffer) fixOffset() {
	if b.Cursor.Offset < b.Cursor.KeepOffset {
		b.Cursor.Offset = b.Cursor.KeepOffset
	}
	if b.Cursor.Offset > len(b.Lines[b.Cursor.Line].Data)-1 {
		b.Cursor.Offset = len(b.Lines[b.Cursor.Line].Data) - 1
	}
	if b.Cursor.Offset < 0 {
		b.Cursor.Offset = 0
	}
}

// Move Cursor left
func (b *Buffer) H(amount int) {
	if b.Cursor.Offset-amount > 0 {
		b.Cursor.Offset -= amount
	} else {
		b.Cursor.Offset = 0
	}
	b.Cursor.KeepOffset = b.Cursor.Offset
}

// Move Cursor down
func (b *Buffer) J(amount int) {
	if b.Cursor.Line+amount < len(b.Lines)-1 {
		b.Cursor.Line += amount
	} else {
		b.Cursor.Line = len(b.Lines) - 1
	}
	b.fixOffset()
}

// Move Cursor up
func (b *Buffer) K(amount int) {
	if b.Cursor.Line-amount > 0 {
		b.Cursor.Line -= amount
	} else {
		b.Cursor.Line = 0
	}
	b.fixOffset()
}

// Move Cursor right
func (b *Buffer) L(amount int) {
	if b.Cursor.Offset+amount < len(b.Lines[b.Cursor.Line].Data) {
		b.Cursor.Offset += amount
	} else {
		b.Cursor.Offset = len(b.Lines[b.Cursor.Line].Data) - 1
		b.fixOffset()
	}
	b.Cursor.KeepOffset = b.Cursor.Offset
}

func (b *Buffer) MoveToFirstLine() {
	b.Cursor.Line = 0
	b.fixOffset()
}

func (b *Buffer) MoveToLastLine() {
	b.Cursor.Line = len(b.Lines) - 1
	b.fixOffset()
}

// Move buffer.Cursor to the first non-space character in the line
func (b *Buffer) MoveToFirst() {
	for i := range len(b.Lines[b.Cursor.Line].Data) {
		if b.Lines[b.Cursor.Line].Data[i] != ' ' {
			b.Cursor.Offset = i
			break
		}
	}
	b.Cursor.KeepOffset = b.Cursor.Offset
}
