package render

import "fmt"

type asciiCode string

func (a asciiCode) str() string {
	return string(a)
}

const (
	reset       asciiCode = "\033[0m"
	symbolColor asciiCode = "\033[90m"

	bold       asciiCode = "\033[1m"
	italic     asciiCode = "\033[3m"
	boldItalic asciiCode = "\033[1m\033[3m"
	underline  asciiCode = "\033[4m"
	stricked   asciiCode = "\033[9m"

	quote       asciiCode = "\033[32m"
	quoteSymbol asciiCode = "\u2503"
	codeLine    asciiCode = "\033[33m"
	header      asciiCode = "\033[94m"
	shield      asciiCode = "\\"
	link        asciiCode = "\033[4;36m"

	listColor asciiCode = "\033[35m"
	tagColor  asciiCode = "\033[35m"
	tagS                = "["
	tagE                = "]"
	listDash  asciiCode = "\u2981"
	boxEmpty  asciiCode = " \u25a1"
	boxField  asciiCode = " \u25a0"
)

type Renderer struct {
	curAttr asciiCode
	l       *lexer
}

func InitReder(w, h int) *Renderer {
	r := &Renderer{}
	// TODO: create a new lexer for code and separate it from the default markdown lexer and renderer
	r.l = NewLexer()
	return r
}

func (r *Renderer) RednerMarkdownLine(line []rune, show bool) (string, int) {
	r.l.input = line
	r.l.position = 0
	r.l.readPosition = 0
	r.l.readChar()

	var data = ""
	var diff = 0

	var i = 0
	for tok := r.l.NextToken(); tok.Type != EOL; tok = r.l.NextToken() {
		switch tok.Type {
		case Header_1, Header_2, Header_3, Header_4, Header_5, Header_6:
			if i == 0 {
				data += r.renderHeader(&tok)
			} else {
				data += string(tok.Literal)
			}
		case ListBoxField:
			if i == 0 {
				data += r.renderBoxField(&tok, show)
			} else {
				data += string(tok.Literal)
			}
		case ListBoxEmpty:
			if i == 0 {
				data += r.renderBoxEmpty(&tok, show)
			} else {
				data += string(tok.Literal)
			}
		case ListDash:
			if i == 0 {
				data += r.renderListDash(&tok, show)
			} else {
				data += string(tok.Literal)
			}
		case Quote:
			if i == 0 {
				data += r.renderQuote(&tok, show)
			} else {
				data += string(tok.Literal)
			}
		case ListNumberB, ListNumberDot:
			if i == 0 {
				data += r.renderListNumber(&tok)
			} else {
				data += string(tok.Literal)
			}
		case Link:
			data += r.renderLink(&tok, show)
		case Image:
			data += r.renderImage(&tok, show)
		case CodeLine:
			data += r.simpleAttrRender(codeLine, string(tok.Literal), show)
			diff += 1
		case TEXT:
			data += r.renderText(&tok)
		case Shield:
			data += r.renderShield(&tok, show)
			diff += 1
		case Tag:
			data += r.renderTag(&tok, show)
			diff -= 1
		case OneStar, OneUnderline:
			data += r.simpleAttrRender(italic, string(tok.Literal), show)
			diff += 1
		case TwoStars, TwoUnderlines:
			data += r.simpleAttrRender(bold, string(tok.Literal), show)
			diff += 2
		case ThreeStars, ThreeUnderlines:
			data += r.simpleAttrRender(boldItalic, string(tok.Literal), show)
			diff += 3
		case Stricked:
			data += r.simpleAttrRender(stricked, string(tok.Literal), show)
			diff += 2
		case WhiteSpace:
			data += " "
		case Symbol:
			data += string(tok.Value)
		}
		i += 1
	}

	data += reset.str()
	r.curAttr = reset
	return data, diff
}

func painAsAttr(symbol string) string {
	sym := paintString(symbolColor, symbol)
	return sym + reset.str()
}

func paintString(ascii asciiCode, str string) string {
	var s = ""
	for _, x := range str {
		s += fmt.Sprintf("%s%c", ascii, x)
	}
	return s
}

func (r *Renderer) renderBoxEmpty(t *Token, show bool) string {
	if show {
		return painAsAttr(string(t.Literal))
	} else {
		return boxEmpty.str()
	}
}

func (r *Renderer) renderBoxField(t *Token, show bool) string {
	if show {
		return painAsAttr(string(t.Literal))
	} else {
		return boxField.str()
	}
}

func (r *Renderer) renderListNumber(t *Token) string {
	var s = ""
	s += paintString(listColor, string(t.Value))
	s += paintString(listColor, string(t.Literal))
	s += reset.str()
	return s
}

func (r *Renderer) renderListDash(t *Token, show bool) string {
	if show {
		return painAsAttr(string(t.Literal))
	} else {
		return listDash.str()
	}
}

func (r *Renderer) renderShield(t *Token, show bool) string {
	var s = ""
	if show {
		s += painAsAttr(string(t.Literal))
	}
	s += string(t.Value)
	return s
}

func (r *Renderer) renderQuote(t *Token, show bool) string {
	var s = ""
	s += quote.str()
	if show {
		s += painAsAttr(string(t.Literal))
	} else {
		s += quoteSymbol.str()
	}
	s += reset.str()
	return s
}

func (r *Renderer) renderText(t *Token) string {
	if r.curAttr != reset {
		return paintString(r.curAttr, string(t.Value))
	}
	return string(t.Value)
}

func (r *Renderer) renderTag(t *Token, show bool) string {
	var s = ""
	if !show {
		s += tagColor.str()
		s += tagS
		s += paintString(tagColor, string(t.Literal))
		s += paintString(tagColor, string(t.Value))
		s += tagE
		s += reset.str()
	} else {
		s += paintString(tagColor, string(t.Literal))
		s += paintString(tagColor, string(t.Value))
	}
	s += reset.str()
	return s
}

func (r *Renderer) renderHeader(t *Token) string {
	var s = ""
	s += header.str()
	s += underline.str()
	r.curAttr = header
	s += string(t.Literal)
	return s
}

func (r *Renderer) renderLink(t *Token, show bool) string {
	if show {
		return link.str() + string(t.Literal) + reset.str()
	} else {
		return link.str() + string(t.Value) + reset.str()
	}
}

func (r *Renderer) renderImage(t *Token, show bool) string {
	if show {
		return link.str() + string(t.Literal) + reset.str()
	} else {
		return link.str() + string(t.Value) + reset.str()
	}
}

func (r *Renderer) simpleAttrRender(mode asciiCode, attr string, show bool) string {
	var s = ""
	if r.curAttr == mode {
		r.curAttr = reset
		if show {
			s += painAsAttr(attr)
		}
		s += r.curAttr.str()
	} else {
		r.curAttr = mode
		s += r.curAttr.str()
		if show {
			s += painAsAttr(attr)
		}
	}
	return s
}
