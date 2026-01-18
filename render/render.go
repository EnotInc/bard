package render

import "strings"

const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	italic    = "\033[3m"
	underline = "\033[7m"
	stricked  = "\033[9m"
	header    = "\033[7;1;38;5;255;48;5;236m"
	startSel  = "\033[100m"
	endSel    = "\033[49m"
)

type Renderer struct {
	curMode string
	w, h    int
	l       *lexer
}

func InitReder(w, h int) *Renderer {
	r := &Renderer{
		w: w,
		h: h,
	}
	return r
}

func (r *Renderer) RednerMarkdownLine(line []rune, isCur bool) string {
	r.l = NewLexer(line)
	var data = ""

	var i = 0
	for tok := r.l.NextToken(); tok.Type != EOL; tok = r.l.NextToken() {
		switch tok.Type {
		case Header_1, Header_2, Header_3, Header_4, Header_5, Header_6:
			if i == 0 {
				data += r.renderHeader(&tok, isCur)
			} else {
				data += string(tok.Literal)
			}
		case TEXT:
			data += string(tok.Literal)
		case OneStar, OneUnderline:
			data += r.renderItalc(&tok, isCur)
		case TwoStars, TwoUnderlines:
			data += r.renderBold(&tok, isCur)
		case Stricked:
			data += r.renderStriked(&tok, isCur)
		case Symbol:
			data += string(tok.Literal)
		}
		i += 1
	}

	data += reset
	r.curMode = reset
	return data
}

func (r *Renderer) renderHeader(t *Token, isCurl bool) string {
	var s = ""
	if !isCurl {
		tabs := r.w / 4
		r.curMode = reset

		s += strings.Repeat(" ", tabs)
	}
	s += header
	//s += bold
	s += string(t.Literal)
	return s
}

func (r *Renderer) renderItalc(t *Token, isCur bool) string {
	var s = ""
	r.changeMode(italic)
	s += r.curMode
	if isCur {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) renderBold(t *Token, show bool) string {
	var s = ""
	//TODO: Figure out why it's now working in windws powershell
	r.changeMode(bold)
	s += r.curMode
	if show {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) renderStriked(t *Token, show bool) string {
	var s = ""
	r.changeMode(stricked)
	s += r.curMode
	if show {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) changeMode(mode string) {
	if r.curMode == mode {
		r.curMode = reset
	} else {
		r.curMode = mode
	}
}
