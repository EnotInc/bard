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
)

func getLogPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".bard"
	}
	return filepath.Join(home, ".bard")
}

func (s *Screen) saveLog(error string) error {
	path := getLogPath()
	logs := filepath.Join(path, ".log")

	file, err := os.OpenFile(logs, os.O_APPEND|os.O_CREATE, 0644)
	theme := config.GetTheme().General
	if err != nil {
		return fmt.Errorf("%s%s%s%s%s", theme.Message, err, "\n\n Error stack:\n", ascii.Reset, string(debug.Stack()))
	}
	defer file.Close()

	log.SetOutput(file)
	log.Print(strings.Repeat("=", 30), "\n\n", error, "\n", string(debug.Stack()), "\n\n")
	return nil
}
