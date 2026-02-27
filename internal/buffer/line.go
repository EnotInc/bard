package buffer

import "slices"

type Line struct {
	Data []rune
}

func (b *Buffer) InsertEmptyLine(lineShift int) {
	index := b.Cursor.Line + lineShift
	newLine := make([]*Line, 0)
	newLine = append(newLine, &Line{Data: []rune("")})
	b.Lines = append(b.Lines[:index], append(newLine, b.Lines[index:]...)...)
}

func (b *Buffer) InsertLineWithData(index int, data []rune) {
	newLine := make([]*Line, 0)
	newLine = append(newLine, &Line{Data: data})
	b.Lines = append(b.Lines[:index], append(newLine, b.Lines[index:]...)...)
}

// Called when the user presses [enter] in the middle of a line. This function shifts data from the right to the new line
func (b *Buffer) InsertLine() {
	index := b.Cursor.Line + 1
	shiftData := b.Lines[b.Cursor.Line].Data[b.Cursor.Offset:]
	b.Lines[b.Cursor.Line].Data = b.Lines[b.Cursor.Line].Data[:b.Cursor.Offset]

	newLine := Line{}
	b.Lines = append(b.Lines[:index], append([]*Line{&newLine}, b.Lines[index:]...)...)
	b.Cursor.Line += 1
	b.Cursor.Offset = 0

	b.Lines[b.Cursor.Line].Data = append(b.Lines[b.Cursor.Line].Data, shiftData...)
}

// Called when the user deletes the 0th character in a line. The line is deleted and data is moved to the line above
func (b *Buffer) DelAndMoveLine() {
	if b.Cursor.Line > 0 {
		shiftData := b.Lines[b.Cursor.Line].Data[b.Cursor.Offset:]
		b.RemoveLine()
		b.Cursor.Line -= 1
		b.Cursor.Offset = len(b.Lines[b.Cursor.Line].Data)
		b.Lines[b.Cursor.Line].Data = append(b.Lines[b.Cursor.Line].Data, shiftData...)
	}
}

func (b *Buffer) DelAndMoveLineAt(startLine int, endLine int, endOffset int) {
	shiftData := b.Lines[endLine].Data[endOffset:]
	b.RemoveLineAt(endLine)
	b.Cursor.Offset = len(b.Lines[startLine].Data)
	b.Lines[startLine].Data = append(b.Lines[startLine].Data, shiftData...)
}

// Delete the whole line
func (b *Buffer) RemoveLine() {
	if len(b.Lines) == 1 {
		b.ClearLine()
		return
	}
	b.Lines = slices.Delete(b.Lines, b.Cursor.Line, b.Cursor.Line+1)
	if b.Cursor.Line >= len(b.Lines) {
		b.Cursor.Line = len(b.Lines)
	}
}

// Delete the whole line at index
func (b *Buffer) RemoveLineAt(lineIndex int) {
	if len(b.Lines) == 1 {
		b.ClearLine()
		return
	}
	b.Lines = slices.Delete(b.Lines, lineIndex, lineIndex+1)
	if b.Cursor.Line >= len(b.Lines) {
		b.K(1)
	}
}

// Set line.data = ""
func (b *Buffer) ClearLine() {
	b.Lines[b.Cursor.Line].Data = []rune{}
}
