package editor

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/EnotInc/Bard/docs/help"
	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/buffer"
	"github.com/EnotInc/Bard/internal/enums"
)

// Used to create hew Buffer in Editor with selected help topic
func (e *Editor) OpenHelp(topic enums.Help) {
	var lines iter.Seq[string]
	switch topic {
	case enums.HelpAbout:
		lines = strings.SplitSeq(help.About, "\n")
	case enums.HelpModes:
		lines = strings.SplitSeq(help.Modes, "\n")
	case enums.HelpNormal:
		lines = strings.SplitSeq(help.Noraml, "\n")
	case enums.HelpCommand:
		lines = strings.SplitSeq(help.Command, "\n")
	case enums.HelpVisual, enums.HelpVisualLine:
		lines = strings.SplitSeq(help.Visual, "\n")
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

// Used to read file data, and write it into current Buffer
func (e *Editor) LoadFile(file string) {
	if f, err := os.Stat(file); err != nil {
		e.CreateFile(file)
		e.tui.ShowHello = true
	} else if f.IsDir() {
		fmt.Printf("'%s' is a dir, not file", file)
		os.Exit(1)
		return
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

// Called when new file is opened or created in Bard
func (e *Editor) CreateFile(fileName string) {
	_, err := os.Create(fileName)
	if err != nil {
		e.tui.Message = fmt.Sprintf("Unable to create file %s", fileName)
	}
}

// saves current Buffer into file
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

func getLogPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".bard"
	}
	return filepath.Join(home, ".bard")
}

func (e *Editor) saveLog(msg any) error {
	path := getLogPath()
	logs := filepath.Join(path, ".log")

	file, err := os.OpenFile(logs, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("%s%s%s%s%s", e.theme.General.Message, msg, "\n\n Error stack:\n", ascii.Reset, string(debug.Stack()))
	}
	defer file.Close()

	log.SetOutput(file)
	log.Print(strings.Repeat("=", 30), "\n\n", msg, "\n", string(debug.Stack()), "\n\n")
	return nil
}
