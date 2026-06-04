package screen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"golang.org/x/term"
)

var global *Screen

type Screen struct {
	oldState *term.State
	tiles    []*tile // let's assume that there is only row layout. I'll figure out this later
	focus    int
	w, h     int
	fdIn     int
	fdOut    int
	status   string
}

func InitScreen() {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil {
		panic(err)
	}

	_w, _h, _ := term.GetSize(_fdOut)
	if _w <= 40 || _h < 15 {
		panic("Unable to run Bard. Window size is too small!")
	}

	s := &Screen{
		oldState: old,
		tiles:    make([]*tile, 0),
		focus:    0,
		w:        0,
		h:        0,
		fdIn:     _fdIn,
		fdOut:    _fdOut,
	}
	global = s
}

// NOTE: I don't rly shure will this work...
func DrawAll() {
	var data strings.Builder
	for _, t := range global.tiles {
		tile := t.GetDiff()
		data.WriteString(tile)
	}
	fmt.Print(data.String())
}

func Run() {
	defer Exit(1)
	fmt.Print(ascii.SaveTerminal, ascii.ClearView, ascii.ClearHistory)
	DrawAll()
	reader := bufio.NewReader(os.Stdin)
	for {
		key, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				key = keys.Ctrl_z
			} else {
				panic(err)
			}
		}

		global.tiles[global.focus].object.Handle(key)
		DrawAll()
	}
}
