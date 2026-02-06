package editor

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func (e *Editor) LoadFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ext := filepath.Ext(file)
	e.isMdFile = (ext == ".md" || ext == ".MD")

	scanner := bufio.NewScanner(f)

	//clearing the list of lines, coz I make one line in InitEditor() func
	e.b.lines = []*line{}
	for scanner.Scan() {
		l := &line{}

		scannedLine := scanner.Text()
		scannedLine = strings.ReplaceAll(scannedLine, "\t", "    ")
		e.b.lines = append(e.b.lines, l)
		e.b.lines[len(e.b.lines)-1].data = []rune(scannedLine)
	}
	if len(e.b.lines) == 0 {
		e.b.lines = append(e.b.lines, &line{})
	}
	e.file = file
}

func (e *Editor) SaveFile() {
	if !(e.file == "" || len(e.file) == 0) {
		var data []byte

		for _, v := range e.b.lines {
			byteLine := []byte(string(v.data))
			data = append(data, byteLine...)
			data = append(data, byte('\n'))
		}

		err := os.WriteFile(e.file, data, 0644)
		if err != nil {
			e.message = err.Error()
		} else {
			e.message = "file saved"
		}
	} else {
		e.message = "file name was not provided"
	}
}
