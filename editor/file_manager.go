package editor

import (
	"bufio"
	"os"
)

func (e *Editor) LoadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)

	//clearing the list of lines, coz I make one line in InitEditor() func
	e.b.lines = []*line{}
	for scanner.Scan() {
		l := &line{}

		scannedLine := scanner.Text()
		e.b.lines = append(e.b.lines, l)
		e.b.lines[len(e.b.lines)-1].data = []rune(scannedLine)
	}
	e.file = file
	return nil
}

func (e *Editor) SaveFile() error {
	if !(e.file == "" || len(e.file) == 0) {
		var data []byte

		for _, v := range e.b.lines {
			byteLine := []byte(string(v.data))
			data = append(data, byteLine...)
			data = append(data, byte('\n'))
		}

		err := os.WriteFile(e.file, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
