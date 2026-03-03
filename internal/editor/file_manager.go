package editor

import (
	"Enot/Bard/docs/help"
	"Enot/Bard/internal/buffer"
	"Enot/Bard/internal/enums"
	"bufio"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

func (e *Editor) OpenHelp(topic enums.Help) {
	e.newBuffer()
	e.b[e.curBuffer].Lines = []*buffer.Line{}
	e.b[e.curBuffer].IsMdFile = true
	e.b[e.curBuffer].Title = string(topic)

	var lines iter.Seq[string]
	switch topic {
	case enums.About:
		lines = strings.SplitSeq(help.About, "\n")
	default:
		e.tui.Message = "Unable to open this help topic"
		return
	}

	for line := range lines {
		l := &buffer.Line{}

		line = strings.ReplaceAll(line, "\t", "    ")
		e.b[e.curBuffer].Lines = append(e.b[e.curBuffer].Lines, l)
		e.b[e.curBuffer].Lines[len(e.b[e.curBuffer].Lines)-1].Data = []rune(line)
	}
}

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

func (e *Editor) CreateFile(fileName string) {
	_, err := os.Create(fileName)
	if err != nil {
		e.tui.Message = fmt.Sprintf("Unable to create file %s", fileName)
	}
}

func (e *Editor) SaveFile() {
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
