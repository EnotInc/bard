package editor

import (
	"Enot/Bard/internal/buffer"
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

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
	e.isMdFile = (ext == ".md" || ext == ".MD")

	scanner := bufio.NewScanner(f)

	//clearing the list of lines, coz I make one line in InitEditor() func
	e.b.Lines = []*buffer.Line{}
	for scanner.Scan() {
		l := &buffer.Line{}

		scannedLine := scanner.Text()
		scannedLine = strings.ReplaceAll(scannedLine, "\t", "    ")
		e.b.Lines = append(e.b.Lines, l)
		e.b.Lines[len(e.b.Lines)-1].Data = []rune(scannedLine)
	}
	if len(e.b.Lines) == 0 {
		e.b.Lines = append(e.b.Lines, &buffer.Line{})
	}
	e.file = file
}

func (e *Editor) CreateFile(fileName string) {
	//TODO: check if fileName is legit
	os.Create(fileName)
}

func (e *Editor) SaveFile() {
	if !(e.file == "" || len(e.file) == 0) {
		var data []byte

		for _, v := range e.b.Lines {
			byteLine := []byte(string(v.Data))
			data = append(data, byteLine...)
			data = append(data, byte('\n'))
		}

		err := os.WriteFile(e.file, data, 0644)
		if err != nil {
			e.tui.Message = err.Error()
		} else {
			ext := filepath.Ext(e.file)
			e.isMdFile = (ext == ".md" || ext == ".MD")

			e.tui.Message = "file saved"
		}
	} else {
		e.tui.Message = "file name was not provided"
	}
}
