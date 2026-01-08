package editor

import (
	"fmt"
	"os"
	"slices"

	"golang.org/x/term"
)

const (
	initialCurOfset  = 2
	initialLineShift = 1
)

const (
	clearView    = "\033[2J"
	clearHistory = "\033[3J"
	moveToStart  = "\033[0H"
	cursorBloc   = "\x1b[0 q"
	cursorLine   = "\x1b[5 q"
)

type Mode string

const (
	normal  = "Normal"
	command = "Command"
	insert  = "Insert"
)

type Buffer struct {
	lines    []string
	curLine  int
	curOfset int
}

func InitBuffer() *Buffer {
	b := &Buffer{
		curLine:  0,
		curOfset: 0,
	}

	b.lines = append(b.lines, "")

	return b
}

type Editor struct {
	oldState   *term.State
	buffer     *Buffer
	curMode    Mode
	curCommand string
	curType    string
	fdIn       int
	w          int
	h          int
}

func InitEditor() *Editor {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil {
		panic(err)
	}
	w, h, _ := term.GetSize(_fdOut)

	b := InitBuffer()

	e := &Editor{
		oldState:   old,
		buffer:     b,
		curMode:    normal,
		curCommand: "",
		fdIn:       _fdIn,
		w:          w,
		h:          h,
	}

	return e
}

func removeKey(s string, index int) string {
	runes := []rune(s)
	index -= 1

	if index < 0 || index >= len(runes) {
		return s
	}

	return string(append(runes[:index], runes[index+1:]...))
}

func (e *Editor) caseNormal(key byte) {
	switch key {
	case 'h':
		e.buffer.H()
	case 'j':
		e.buffer.J()
	case 'k':
		e.buffer.K()
	case 'l':
		e.buffer.L()
	case 'i':
		e.curMode = insert
	case 'a':
		if e.buffer.curOfset < len(e.buffer.lines[e.buffer.curLine]) {
			e.buffer.curOfset += 1
		}
		e.curMode = insert
	case ':':
		e.curMode = command
	case 'o':
		//p.buffer.lines = append(p.buffer.lines, "")
		e.insertNewLine(1)
		e.buffer.curLine += 1
		e.buffer.curOfset = 0
		e.curMode = insert
	case 'O':
		if e.buffer.curLine > 0 {
			e.insertNewLine(0)
			e.buffer.curOfset = 0
			e.curMode = insert
		}
	}
}

func (e *Editor) caseCommand(key byte) {
	var lastCursorLine int
	var lastCursorOfset int
	if len(e.curCommand) == 0 {
		lastCursorLine = e.buffer.curLine
		lastCursorOfset = e.buffer.curOfset
	}
	e.buffer.curLine = e.h - 1
	e.buffer.curOfset = len(e.curCommand)

	Draw(e)

	switch key {
	case '\033', '\003':
		e.curCommand = ""
		e.buffer.curLine = lastCursorLine
		e.buffer.curOfset = lastCursorOfset
		e.curMode = normal
	case '\127', '\x7f':
		e.curCommand = removeKey(e.curCommand, len(e.curCommand))
		e.buffer.curOfset -= 2 //NOTE: magic number 2, it's just exists

		if len(e.curCommand) == 0 {
			e.buffer.curLine = lastCursorLine
			e.buffer.curOfset = lastCursorOfset
			e.curMode = normal
		}
	case '\013', '\r', '\n':
		success := e.execCommand()
		if !success {
			e.buffer.curLine = lastCursorLine
			e.buffer.curOfset = lastCursorOfset
			e.curCommand = ""
			e.curMode = normal
		}
	default:
		e.curCommand += string(key)
	}
}

func (e *Editor) execCommand() bool {
	switch e.curCommand {
	case "q":
		//TODO: figure out how to propertly resotre the old state
		fmt.Printf("%s%s%s", clearView, clearHistory, moveToStart)
		term.Restore(e.fdIn, e.oldState)
		fmt.Fprintf(os.Stdout, cursorBloc)
		os.Exit(0)
		return true
	default:
		return false
	}
}

func (e *Editor) caseInsert(key byte) {
	switch key {
	case '\013', '\r', '\n':
		e.insertLine()
		e.buffer.curLine += 1
		e.buffer.curOfset = 0
	case '\033', '\003':
		e.curMode = normal
	case '\127', '\x7f':
		if e.buffer.curOfset > 0 {
			e.buffer.lines[e.buffer.curLine] = removeKey(e.buffer.lines[e.buffer.curLine], e.buffer.curOfset)
			e.buffer.curOfset -= 1
		} else {
			if e.buffer.curLine > 0 {
				shiftData := e.removeLine()
				e.buffer.curLine -= 1
				e.buffer.curOfset = len(e.buffer.lines[e.buffer.curLine])
				e.buffer.lines[e.buffer.curLine] += shiftData
			}
		}
	default:
		e.buffer.lines[e.buffer.curLine] = e.insertKey(key)
		if key == '\t' {
			e.buffer.curOfset += 4
		} else {
			e.buffer.curOfset += 1
		}
	}
}

func (e *Editor) insertKey(key byte) string {
	curLine := e.buffer.lines[e.buffer.curLine]

	return curLine[:e.buffer.curOfset] + string(key) + curLine[e.buffer.curOfset:]
}

func (e *Editor) insertNewLine(lineShift int) {
	//NOTE: 1 - line below, 0 - line above
	lineGap := e.buffer.curLine + lineShift

	newLines := append(e.buffer.lines[:lineGap], append([]string{""}, e.buffer.lines[lineGap:]...)...)
	e.buffer.lines = newLines
}

func (e *Editor) insertLine() {
	lineShift := 1
	lineGap := e.buffer.curLine + lineShift

	var newLines []string
	var shiftData string = ""
	if e.buffer.curOfset < len(e.buffer.lines[e.buffer.curLine]) {
		shiftData = e.buffer.lines[e.buffer.curLine][e.buffer.curOfset:]
		e.buffer.lines[e.buffer.curLine] = e.buffer.lines[e.buffer.curLine][:e.buffer.curOfset]
	}
	//TODO: figure this out...
	newLines = append(e.buffer.lines[:lineGap], append([]string{shiftData}, e.buffer.lines[lineGap:]...)...)
	e.buffer.lines = newLines
}

func (e *Editor) removeLine() string {
	var shiftData string = ""
	if len(e.buffer.lines[e.buffer.curLine]) > 0 {
		shiftData = e.buffer.lines[e.buffer.curLine]
	}
	index := e.buffer.curLine
	e.buffer.lines = slices.Delete(e.buffer.lines, index, index+1)
	return shiftData
}

func (b *Buffer) H() {
	if b.curOfset > 0 {
		b.curOfset -= 1
	} else if b.curLine > 0 {
		b.curLine -= 1
		b.curOfset = len(b.lines[b.curLine])
	}
}
func (b *Buffer) J() {
	if b.curLine < len(b.lines)-1 {
		b.curLine += 1
		if b.curOfset > len(b.lines[b.curLine]) {
			b.curOfset = len(b.lines[b.curLine])
		}
	}
}
func (b *Buffer) K() {
	if b.curLine > 0 {
		b.curLine -= 1
		if b.curOfset > len(b.lines[b.curLine]) {
			b.curOfset = len(b.lines[b.curLine])
		}
	}

}
func (b *Buffer) L() {
	if b.curOfset < len(b.lines[b.curLine]) {
		b.curOfset += 1
	} else if b.curLine < len(b.lines)-1 {
		b.curLine += 1
		b.curOfset = 0
	}
}

func (e *Editor) Run() {
	var buf [1]byte
	Draw(e)
	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil {
			panic(err)
		}

		key := buf[0]
		switch e.curMode {
		case normal:
			e.curType = cursorBloc
			e.caseNormal(key)
		case command:
			e.caseCommand(key)
		case insert:
			e.curType = cursorLine
			e.caseInsert(key)
		}

		Draw(e)
	}
}
