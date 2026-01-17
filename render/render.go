package render

import (
	"fmt"
)

const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	italic    = "\033[3m"
	underline = "\033[7m"
	stricked  = "\033[9m"
	startSel  = "\033[100m"
	endSel    = "\033[49m"
)

type Renderer struct {
	curMod      string
	width       int
	HStartLine  int
	HStartOfset int
	HEndLine    int
	HEndOfset   int
}

func InitReder(w int) *Renderer {
	r := &Renderer{
		curMod:      reset,
		width:       w,
		HStartLine:  -1,
		HStartOfset: -1,
		HEndLine:    -1,
		HEndOfset:   -1,
	}

	return r
}

func (r *Renderer) RenderVisualLine(line []rune, lineIndex int, isCurLine bool) string {
	l := ""

	startO := r.HStartOfset
	endO := r.HEndOfset
	startL := r.HStartLine
	endL := r.HEndLine

	if endO+(endL*r.width) < startO+(startL*r.width) {
		startO, endO = endO, startO
		startL, endL = endL, startL
	}
	// if endL < startL {
	// }
	if startL == endL && endL == lineIndex {
		l += string(line[:startO])
		l += startSel
		l += string(line[startO:endO])
		l += endSel
		l += string(line[endO:])
	} else if lineIndex == startL {
		l += string(line[:startO])
		l += startSel
		l += string(line[startO:])
	} else if lineIndex == endL {
		l += startSel
		l += string(line[:endO])
		l += endSel
		l += string(line[endO:])
	} else if startL < lineIndex && lineIndex < endL {
		l = fmt.Sprintf("%s%s", startSel, string(line))
	} else {
		l = fmt.Sprintf("%s", string(line))
	}
	l += fmt.Sprintf("%s", endSel)
	return l
}

func (r *Renderer) RednerMarkDownLine(line []rune, isCurLine bool) string {
	//TODO; meke markdown render
	data := string(line)
	return data
}
