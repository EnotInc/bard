package screen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"golang.org/x/term"
)

func Exit(code int) {
	config.Save()

	fmt.Print(ascii.ClearView, ascii.ClearHistory, ascii.MoveToStart, ascii.CursorReset, ascii.ResetTerminal, ascii.ResetCursor)
	term.Restore(global.fdIn, global.oldState)
	if r := recover(); r != nil {
		err := global.saveLog(r)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Bard stopped with error. More information you can find in '~/.bard/.log' file")
		}
	}
	os.Exit(code)
}

func getLogPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".bard"
	}
	return filepath.Join(home, ".bard")
}

func (s *Screen) saveLog(err any) error {
	path := getLogPath()
	logs := filepath.Join(path, ".log")

	file, err := os.OpenFile(logs, os.O_APPEND|os.O_CREATE, 0644)
	theme := config.GetTheme().General
	if err != nil {
		return fmt.Errorf("%s%s%s%s%s", theme.Message, err, "\n\n Error stack:\n", ascii.Reset, string(debug.Stack()))
	}
	defer file.Close()

	log.SetOutput(file)
	log.Print(strings.Repeat("=", 30), "\n\n", err, "\n", string(debug.Stack()), "\n\n")
	return nil
}
