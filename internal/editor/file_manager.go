package editor

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/EnotInc/Bard/docs/help"
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
)

// About OpenHelp()
// Used to create hew [Buffer] in [Editor] with selected help topic
func (e *Editor) OpenHelp(topic enums.Help) {
	var lines iter.Seq[string]
	switch topic {
	case enums.About:
		lines = strings.SplitSeq(help.About, "\n")
	case enums.Modes:
		lines = strings.SplitSeq(help.Modes, "\n")
	case enums.Normal:
		lines = strings.SplitSeq(help.Noraml, "\n")
	case enums.Command:
		lines = strings.SplitSeq(help.Command, "\n")
	default:
		e.tui.Message = "Unknown topic"
		return
	}

	e.newBuffer()
	e.b[e.curBuffer].Lines = []*buffer.Line{}
	e.b[e.curBuffer].IsReadOnly = true
	e.b[e.curBuffer].IsMdFile = true
	e.b[e.curBuffer].Title = string(topic)

	for line := range lines {
		l := &buffer.Line{}

		line = strings.ReplaceAll(line, "\t", "    ")
		e.b[e.curBuffer].Lines = append(e.b[e.curBuffer].Lines, l)
		e.b[e.curBuffer].Lines[len(e.b[e.curBuffer].Lines)-1].Data = []rune(line)
	}
}

// About LoadFile()
// Used to read file data, and write it into current [Buffer]
func (e *Editor) LoadFile(file string) {
	if _, err := os.Stat(file); err != nil {
		e.CreateFile(file)
		e.tui.ShowHello = true
	}

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ext := filepath.Ext(file)
	e.b[e.curBuffer].IsMdFile = (ext == ".md" || ext == ".MD")

	scanner := bufio.NewScanner(f)

	//clearing the list of lines, coz I make one line in InitEditor() func
	e.b[e.curBuffer].Lines = []*buffer.Line{}
	for scanner.Scan() {
		l := &buffer.Line{}

		scannedLine := scanner.Text()
		scannedLine = strings.ReplaceAll(scannedLine, "\t", "    ")
		e.b[e.curBuffer].Lines = append(e.b[e.curBuffer].Lines, l)
		e.b[e.curBuffer].Lines[len(e.b[e.curBuffer].Lines)-1].Data = []rune(scannedLine)
	}
	if len(e.b[e.curBuffer].Lines) == 0 {
		e.b[e.curBuffer].Lines = append(e.b[e.curBuffer].Lines, &buffer.Line{})
	}
	e.b[e.curBuffer].Title = file
}

// About CreateFile()
// Called when new file is opened or created in Bard
func (e *Editor) CreateFile(fileName string) {
	_, err := os.Create(fileName)
	if err != nil {
		e.tui.Message = fmt.Sprintf("Unable to create file %s", fileName)
	}
}

// About SaveFile()
// saves current [Buffer] into file
func (e *Editor) SaveFile() {
	if !e.b[e.curBuffer].IsReadOnly {
		if !(e.b[e.curBuffer].Title == "" || len(e.b[e.curBuffer].Title) == 0) {
			var data []byte

			for _, v := range e.b[e.curBuffer].Lines {
				byteLine := []byte(string(v.Data))
				data = append(data, byteLine...)
				data = append(data, byte('\n'))
			}

			err := os.WriteFile(e.b[e.curBuffer].Title, data, 0644)
			if err != nil {
				e.tui.Message = err.Error()
			} else {
				ext := filepath.Ext(e.b[e.curBuffer].Title)
				e.b[e.curBuffer].IsMdFile = (ext == ".md" || ext == ".MD")

				e.tui.Message = "file saved"
			}
		} else {
			e.tui.Message = "file name was not provided"
		}
	}
}
