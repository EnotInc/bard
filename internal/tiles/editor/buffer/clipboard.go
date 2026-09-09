package buffer

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
)

// TODO: later implement pasteFromClipboard
func (b *Buffer) copyToClipboard(isVisualLine bool) {
	if !config.GetConfig().ClipBoard {
		return
	}

	var s strings.Builder
	for _, l := range b.Copies {
		s.WriteString(string(l.data))
		if l.isEnd || isVisualLine {
			s.WriteString("\n")
		}
	}
	enc := base64.StdEncoding.EncodeToString([]byte(s.String()))
	fmt.Printf("\x1b]52;c;%s\x07", enc)
}
