package render

import (
	"strings"
)

const (
	reset       = "\033[0m"
	resetColor  = "\033[39m"
	bold        = "\033[1m"
	italic      = "\033[3m"
	underline   = "\033[7m" // so this doesn't work, ig, and now I use this for headers coz it give me a cool bg highlight
	stricked    = "\033[9m"
	symbolColor = "\033[90m"
	quote       = "\033[32m"
	startSel    = "\033[100m"
	endSel      = "\033[49m"
	boldItalic  = "\033[1m\033[3m"

	code   = "\033[48;5;236m"
	header = "\033[7;1;38;5;255;48;5;236m"
	tag    = "\033[48;5;60m\033[38;5;219m"

	tagS     = "["
	tagE     = "]"
	shield   = "\\"
	listDash = "\u2981"
	boxEmpty = "\u2610 "
	boxField = "\u2612 "
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
		case ListBoxField:
			if i == 0 {
				data += r.renderBoxField(&tok, isCur)
			} else {
				data += string(tok.Literal)
			}
		case ListBoxEmpty:
			if i == 0 {
				data += r.renderBoxEmpty(&tok, isCur)
			} else {
				data += string(tok.Literal)
			}
		case ListDash:
			if i == 0 {
				data += r.renderListDash(&tok, isCur)
			} else {
				data += string(tok.Literal)
			}
		case Quote:
			if i == 0 {
				data += r.renderQuote(&tok, isCur)
			} else {
				data += string(tok.Literal)
			}
		case ListNumberB, ListNumberDot:
			if i == 0 {
				data += r.renderListNumber(&tok)
			} else {
				data += string(tok.Literal)
			}
		case CodeLine:
			data += r.renderCodeLine(&tok, isCur)
		case TEXT:
			data += r.renderText(&tok)
		case Shield:
			data += r.renderShield(isCur)
		case Tag:
			data += r.renderTag(&tok, isCur)
		case OneStar, OneUnderline:
			data += r.renderItalc(&tok, isCur)
		case TwoStars, TwoUnderlines:
			data += r.renderBold(&tok, isCur)
		case ThreeStars, ThreeUnderlines:
			data += r.rednerBoldItalic(&tok, isCur)
		case Stricked:
			data += r.renderStriked(&tok, isCur)
		case WhiteSpace:
			data += " "
		case Symbol, Unknow:
			data += string(tok.Literal)
		}
		i += 1
	}

	data += reset
	r.curAttr = reset
	return data
}

func colorise(symbol string) string {
	return symbolColor + symbol + resetColor
}

func (r *Renderer) renderBoxEmpty(t *Token, isCur bool) string {
	var s = ""
	if isCur {
		s += colorise(string(t.Literal))
	} else {
		s += boxEmpty
	}
	return s
}

func (r *Renderer) renderBoxField(t *Token, isCur bool) string {
	var s = ""
	if isCur {
		s += colorise(string(t.Literal))
	} else {
		s += boxField
	}
	return s
}

func (r *Renderer) renderListNumber(t *Token) string {
	var s = ""
	s += colorise(string(t.Literal))
	return s
}

func (r *Renderer) renderListDash(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		if isCur {
			s += colorise(string(t.Literal))
		} else {
			s += listDash
		}
	} else {
		s += string(t.Literal)
		r.shielded = false
	}
	return s
}

func (r *Renderer) renderShield(isCur bool) string {
	var s = ""
	if !r.shielded {
		r.shielded = true
		if isCur {
			s += colorise("\\")
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
			s += colorise(string(t.Literal))
		} else {
			s += "\u2503"
		}
		s += reset
	} else {
		s += string(t.Literal)
		r.shielded = false
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
			s += colorise(string(t.Literal))
			r.shielded = false
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
		r.shielded = false
	}
	return s
}

func (r *Renderer) renderItalc(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(italic)
		s += r.curAttr
		if isCur {
			s += colorise(string(t.Literal))
		}
	} else {
		s += string(t.Literal)
		r.shielded = false
	}

	return s
}

func (r *Renderer) renderBold(t *Token, show bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(bold)
		s += r.curAttr
		if show {
			s += colorise(string(t.Literal))
		}
	} else {
		s += string(t.Literal)
		r.shielded = false
	}

	return s
}

func (r *Renderer) rednerBoldItalic(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(boldItalic)
		s += r.curAttr
		if isCur {
			s += colorise(string(t.Literal))
		}
	} else {
		s += string(t.Literal)
		r.shielded = false
	}

	return s
}

func (r *Renderer) renderStriked(t *Token, show bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(stricked)
		s += r.curAttr
		if show {
			s += colorise(string(t.Literal))
		}
	} else {
		s += string(t.Literal)
		r.shielded = false
	}

	return s
}

func (r *Renderer) renderCodeLine(t *Token, isCur bool) string {
	var s = ""
	if !r.shielded {
		r.changeMode(code)
		s += r.curAttr
		if isCur {
			s += colorise(string(t.Literal))
		}
	} else {
		s += string(t.Literal)
		r.shielded = false
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
