package render

import (
	"strings"
)

const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	italic    = "\033[3m"
	underline = "\033[7m" // so this doesn't work, ig, and now I use this for headers coz it give me a cool bg highlight
	stricked  = "\033[9m"
	shield    = "\033[90m"
	code      = "\033[48;5;236m"
	header    = "\033[7;1;38;5;255;48;5;236m"
	tag       = "\033[48;5;60m\033[38;5;219m"
	tagS      = "["
	tagE      = "]"
	quote     = "\033[32m"
	startSel  = "\033[100m"
	endSel    = "\033[49m"
)

type Renderer struct {
	curAttr  string
	w, h     int
	shielded bool
	l        *lexer
}

func InitReder(w, h int) *Renderer {
	r := &Renderer{
		w:        w,
		h:        h,
		shielded: false,
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
			r.shielded = false
		case Quote:
			if i == 0 {
				data += r.renderQuote(&tok, isCur)
			} else {
				data += string(tok.Literal)
			}
			r.shielded = false
		case CodeLine:
			data += r.renderCodeLine(&tok, isCur)
		case TEXT:
			data += r.renderText(&tok)
			r.shielded = false
		case Shield:
			data += r.renderShield(isCur)
		case Tag:
			data += r.renderTag(&tok, isCur)
			r.shielded = false
		case OneStar, OneUnderline:
			data += r.renderItalc(&tok, isCur)
			r.shielded = false
		case TwoStars, TwoUnderlines:
			data += r.renderBold(&tok, isCur)
			r.shielded = false
		case Stricked:
			data += r.renderStriked(&tok, isCur)
			r.shielded = false
		case WhiteSpace:
			data += " "
			r.shielded = false
		case Symbol, Unknow:
			data += reset
			data += string(tok.Literal)
			r.shielded = false
		}
		i += 1
	}

	data += reset
	r.curAttr = reset
	return data
}

func (r *Renderer) renderShield(isCur bool) string {
	var s = ""
	if !r.shielded {
		r.shielded = true
		if isCur {
			s += shield
			s += "\\"
			s += reset
		}
	} else {
		s += "\\"
		r.shielded = false
	}
	return s
}

func (r *Renderer) renderQuote(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		s += quote
		if isCur {
			s += string(t.Literal)
		} else {
			s += "\u2503"
		}
		s += reset
	} else {
		s += string(t.Literal)
	}
	return s
}

func (r *Renderer) renderText(t *Token) string {
	return string(t.Literal)
}

func (r *Renderer) renderTag(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		s += tag
		if !isCur {
			s += tagS
			s += string(t.Literal)
			s += tagE
		} else {
			s += string(t.Literal)
		}
		s += reset
	} else {
		s += string(t.Literal)
	}
	return s
}

func (r *Renderer) renderHeader(t *Token, isCurl bool) string {
	var s = ""
	if !r.shielded {
		if !isCurl {
			tabs := (r.w - len(r.l.input)) / 2
			r.curAttr = reset

			s += strings.Repeat(" ", tabs)
		}
		//s += header
		s += underline
		s += string(t.Literal)
	} else {
		s += string(t.Literal)
	}
	return s
}

func (r *Renderer) renderItalc(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(italic)
		s += r.curAttr
		if isCur {
			s += string(t.Literal)
		}
	} else {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) renderBold(t *Token, show bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(bold)
		s += r.curAttr
		if show {
			s += string(t.Literal)
		}
	} else {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) renderStriked(t *Token, show bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(stricked)
		s += r.curAttr
		if show {
			s += string(t.Literal)
		}
	} else {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) renderCodeLine(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(code)
		s += r.curAttr
		if isCur {
			s += string(t.Literal)
		}
	} else {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) changeMode(mode string) {
	if r.curAttr == mode {
		r.curAttr = reset
	} else {
		r.curAttr = mode
	}
}
