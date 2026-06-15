package editor

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/EnotInc/Bard/docs/help"
	"github.com/EnotInc/Bard/internal/editor/buffer"
	"github.com/EnotInc/Bard/internal/enums/buffers"
	h "github.com/EnotInc/Bard/internal/enums/help"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
	"github.com/EnotInc/Bard/internal/screen"
)

// Used to create hew Buffer in Editor with selected help topic
func (e *Editor) OpenHelp(topic h.Topic) {
	var lines iter.Seq[string]
	switch topic {
	case h.About:
		lines = strings.SplitSeq(help.About, "\n")
	case h.Modes:
		lines = strings.SplitSeq(help.Modes, "\n")
	case h.Normal:
		lines = strings.SplitSeq(help.Normal, "\n")
	case h.Command:
		lines = strings.SplitSeq(help.Command, "\n")
	case h.Visual, h.VisualLine:
		lines = strings.SplitSeq(help.Visual, "\n")
	case h.Config:
		lines = strings.SplitSeq(help.Config, "\n")
	case h.Explorer:
		lines = strings.SplitSeq(help.Explorer, "\n")
	default:
		e.tui.Message = "Unknown topic"
		return
	}

	e.newBuffer()
	e.b[e.curBuffer].Lines = []*buffer.Line{}
	e.b[e.curBuffer].IsReadOnly = true
	e.b[e.curBuffer].Type = buffers.Markdown
	e.b[e.curBuffer].Title = string(topic)

	for line := range lines {
		l := &buffer.Line{}

		e.b[e.curBuffer].Lines = append(e.b[e.curBuffer].Lines, l)
		e.b[e.curBuffer].Lines[len(e.b[e.curBuffer].Lines)-1].Data = []rune(line)
	}
}

func (e *Editor) StartHelp() {
	e.OpenHelp(h.About)
	e.delBuffer(0)
}

// Used to read file data, and write it into current Buffer
func (e *Editor) LoadFile(file string) {
	if f, err := os.Stat(file); err != nil {
		e.CreateFile(file)
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

	e.setBufferType(file)

	scanner := bufio.NewScanner(f)

	if scanner.Err() != nil {
		screen.Exit(1)
	}

	//clearing the list of lines, coz I make one line in InitEditor() func
	e.b[e.curBuffer].Lines = []*buffer.Line{}
	for scanner.Scan() {
		l := &buffer.Line{}

		scannedLine := scanner.Text()
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
	f, err := os.Create(fileName)
	if err != nil {
		e.tui.Message = fmt.Sprintf("Unable to create file %s", fileName)
	}
	defer f.Close()
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
				//ext := filepath.Ext(e.b[e.curBuffer].Title)
				//e.b[e.curBuffer].IsMdFile = (ext == ".md" || ext == ".MD")
				e.setBufferType(e.b[e.curBuffer].Title)

				e.tui.Message = "file saved"
			}
		} else {
			e.tui.Message = "file name was not provided"
		}
	}
}

func (e *Editor) OpenFileCallback(file string) {
	if len(e.b) == 1 && e.b[0].Title == "" && len(e.b[0].Lines) == 1 {
		e.LoadFile(file)
		return
	}

	for i, b := range e.b {
		title := strings.TrimPrefix(b.Title, ".\\")
		if title == file {
			e.SwitchToTab(i)
			return
		}
	}

	e.newBuffer()
	e.LoadFile(file)
}

func (e *Editor) RemoveFileCallback(file string) {
	e.curMode = mode.Command
	e.cmd.command = fmt.Sprintf("del %s", file)
}

func (e *Editor) ChangeModeCallback(mode mode.Mode) {
	e.curMode = mode
}

func (e *Editor) setBufferType(file string) {
	ext := filepath.Ext(file)
	ext = strings.TrimPrefix(ext, ".")
	ext = strings.ToLower(ext)
	if ext == "md" || ext == "MD" {
		e.b[e.curBuffer].Type = buffers.Markdown
	} else if ok := buffers.CodeExt[ext]; ok {
		e.b[e.curBuffer].Type = buffers.Code
	} else {
		e.b[e.curBuffer].Type = buffers.Other
	}
}
