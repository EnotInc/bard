package screen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"golang.org/x/term"
)

var global *Screen

func Get() *Screen {
	return global
}

func SendCall(c calls.Call) {
	global.call = c
}

type Screen struct {
	oldState *term.State
	tiles    []*tile // let's assume that there is only row layout. I'll figure out this later
	call     calls.Call
	focus    int
	w, h     int
	fdIn     int
	fdOut    int
	status   func(withBorder bool) string
}

func W() int {
	return global.w
}

func H() int {
	return global.h
}

func InitScreen() {
	_fdIn := int(os.Stdin.Fd())
	_fdOut := int(os.Stdout.Fd())

	old, err := term.MakeRaw(_fdIn)
	if err != nil && 1 == 2 {
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
		w:        _w,
		h:        _h,
		fdIn:     _fdIn,
		fdOut:    _fdOut,
	}
	global = s
}

func AddTile(t *tile) {
	global.tiles = append(global.tiles, t)
}

func ShiftFocus() {
	if global.focus == len(global.tiles)-1 {
		global.focus = 0
		return
	}

	global.focus = len(global.tiles) - 1
}

func SetFocus(index int) {
	if index < 0 {
		global.focus = 0
		return
	}
	if index > len(global.tiles)-1 {
		global.focus = len(global.tiles) - 1
		return
	}

	global.focus = index
}

func SetStatusBar(builder func(withBorder bool) string) {
	global.status = builder
}

// NOTE: I don't rly shure will this work...
func DrawAll() {
	handleCalls()

	var data strings.Builder

	for _, t := range global.tiles {
		t.object.PreDraw()
		tile := t.GetDiff()
		data.WriteString(tile)
	}

	f_tile := global.tiles[global.focus]

	ofset := 0
	if f_tile.border {
		ofset = 1
	}

	status := global.status(f_tile.border)
	data.WriteString(status)

	cX, cY := f_tile.object.GetCursor(f_tile.border)
	cX += f_tile.x + ofset
	cY += f_tile.y + ofset

	fmt.Fprintf(&data, "\033[%d;%dH", cY, cX)
	data.WriteString(string(ascii.ShowCursor))

	fmt.Print(data.String())
}

func handleCalls() {
	switch global.call {
	case calls.PurgeCache:
		for i := range global.tiles {
			global.tiles[i].hash = make(map[int]uint32, 0)
		}
	case calls.Rezise:
	}
	global.call = calls.None
}

func Run() {
	defer Exit(1)
	fmt.Print(ascii.SaveTerminal, ascii.ClearView, ascii.ClearHistory)
	DrawAll()
	reader := bufio.NewReader(os.Stdin)
	for {
		// TODO: read as buffer. Add ascii escape sequances parser
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
