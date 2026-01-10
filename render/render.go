package render

import "fmt"

const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	italic    = "\033[3m"
	underline = "\033[7m"
	stricked  = "\033[9m"
)

type Renderer struct {
	curMod string
}

func InitReder() *Renderer {
	r := &Renderer{
		curMod: reset,
	}

	return r
}

func (r *Renderer) RednerLine(line []rune, isCur bool) string {
	var data = ""
	for i := range len(line) {
		cur := line[i]
		if i == 0 && cur == '#' {
			data += fmt.Sprintf("%s%s", bold, string(line))
			break
		}
		switch cur {
		case '_':
			data += r.renderChar(cur, italic, isCur)
		case '-':
			data += r.renderChar(cur, stricked, isCur)
		default:
			data += fmt.Sprintf("%c", cur)
		}
	}
	data += reset
	r.curMod = reset
	return data
}

func (r *Renderer) renderChar(ch rune, mod string, isCur bool) string {
	var data = ""
	switch r.curMod {
	case reset:
		r.curMod = mod
	case mod:
		data += fmt.Sprintf("%s", r.curMod)
		r.curMod = reset
	}
	if isCur {
		data += fmt.Sprintf("%c", ch)
	}
	data += fmt.Sprintf("%s", r.curMod)

	return data
}
