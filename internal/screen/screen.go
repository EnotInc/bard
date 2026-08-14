package screen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/enums/calls"
	"github.com/EnotInc/Bard/internal/enums/keys"
	"github.com/EnotInc/Bard/internal/enums/popups"
	"golang.org/x/term"
)

var global *Screen

func SendCall(c calls.Call) {
	global.call = c
}

type Screen struct {
	resize   chan bool
	oldState *term.State
	status   func(withBorder bool) string
	popups   map[popups.PopUp]*popup
	root     []rune
	tiles    []*tile
	hiden    tile
	call     calls.Call
	focus    int
	popup    popups.PopUp
	w        int
	h        int
	fdIn     int
	fdOut    int
}

func W() int {
	return global.w
}

func H() int {
	return global.h
}

func Root() []rune {
	return global.root
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
		resize:   make(chan bool, 1),
		oldState: old,
		popups:   make(map[popups.PopUp]*popup, 0),
		tiles:    make([]*tile, 0),
		focus:    0,
		popup:    popups.Close,
		w:        _w,
		h:        _h,
		fdIn:     _fdIn,
		fdOut:    _fdOut,
		root:     []rune(enums.DefaultRoot),
	}
	global = s
}

func AddPopup(p *popup, t popups.PopUp) {
	global.popups[t] = p
}

func AddTile(t *tile) {
	global.tiles = append(global.tiles, t)
}

func ShiftFocus() {
	if global.focus == len(global.tiles)-1 {
		global.focus = 0
	} else {
		global.focus += 1
	}
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

func drawPopUp() {
	p := global.popups[global.popup]

	p.object.PreDraw()
	p.Draw()

	var data strings.Builder
	cX, cY, _ := p.object.GetCursor(false)
	fmt.Fprintf(&data, "\033[%d;%dH", cY, cX)
	data.WriteString(string(ascii.HideCursor))

	fmt.Print(data.String())
}

func DrawAll() {
	handleCalls()
	if global.popup != popups.Close { // render popup window
		drawPopUp()
		return
	}

	var data strings.Builder
	var tilesOfset int = 0
	var focusedOfset = 0

	for i, t := range global.tiles {
		focused := i == global.focus
		t.object.PreDraw()
		tile := t.GetDiff(tilesOfset, focused)
		data.WriteString(tile)

		if focused {
			focusedOfset = tilesOfset
		}
		tilesOfset += t.w
	}

	f_tile := global.tiles[global.focus]

	border := config.GetConfig().ShowBorder

	offset := 0
	if border {
		offset = 1
	}

	status := global.status(border)
	fmt.Fprintf(&data, "\033[%d;1H", global.h)
	data.WriteString(status)

	cX, cY, cursor := f_tile.object.GetCursor(border)
	cX += offset + focusedOfset
	cY += offset

	fmt.Fprintf(&data, "\033[%d;%dH", cY, cX)
	data.WriteString(string(cursor))
	data.WriteString(string(ascii.ShowCursor))

	fmt.Print(data.String())
}

func handleCalls() {
	switch global.call {
	case calls.OpenThemes:
		global.popup = popups.Themes

	case calls.ClosePopups:
		global.popup = popups.Close
		for i := range global.tiles {
			global.tiles[i].hash = make(map[int]uint32, 0)
		}

	case calls.PurgeCache:
		for i := range global.tiles {
			global.tiles[i].hash = make(map[int]uint32, 0)
		}
	case calls.OpenFile, calls.DelFile:
		ShiftFocus()
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

		if global.popup != popups.Close {
			global.popups[global.popup].object.Handle(key)
			DrawAll()
			continue
		}

		switch key {
		case keys.Ctrl_o:
			ShiftFocus()
		case keys.Ctrl_j:
			HideTile()
		default:
			global.tiles[global.focus].object.Handle(key)
		}
		DrawAll()
	}
}

func HideTile() {
	global.call = calls.PurgeCache
	if len(global.tiles) == 2 {
		global.hiden = *global.tiles[1]
		global.tiles = global.tiles[:len(global.tiles)-1]
		SetFocus(0)

		ed := global.tiles[0]
		ed.w = global.w
		ed.object.Resize(ed.w, ed.h)
	} else {
		t := global.hiden
		AddTile(&t)
		SetFocus(1)
		global.hiden = tile{}

		ed := global.tiles[0]
		ed.w = global.w - t.w
		ed.object.Resize(ed.w, ed.h)
	}
	global.call = calls.PurgeCache
}

func SetRoot(root string) {
	global.root = []rune(root)
}
