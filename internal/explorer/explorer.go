package explorer

import (
	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	mode "github.com/EnotInc/Bard/internal/enums/mode"
)

type action int

const (
	_ action = iota
	none
	changing
	creating
	deleting
)

type Explorer struct {
	cursor     *cursor
	visible    *cursor
	openFile   func(file string)
	delFile    func(file string)
	rename     func(old, new string)
	changeMode func(mode mode.Mode)
	path       []rune
	entries    []entry
	buffer     entry
	w          int
	h          int
	yScroll    int
	action     action
	update     bool
	showDot    bool
}

func InitExplorer(open, del func(file string), ren func(old, new string), change func(mode mode.Mode), w, h int) *Explorer {

	c := initCursor()
	v := initCursor()
	cfg := config.GetConfig()
	ex := &Explorer{
		path:       []rune(enums.DefaultRoot),
		w:          w,
		h:          h,
		cursor:     c,
		visible:    v,
		openFile:   open,
		delFile:    del,
		changeMode: change,
		rename:     ren,
		action:     none,
		update:     true,
		showDot:    cfg.ShowDot,
	}
	ex.scanEntries()
	ex.scroll()

	return ex
}

func (ex *Explorer) SetPath(path string) {
	ex.path = []rune(path)
}
