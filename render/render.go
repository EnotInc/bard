package render

import "fmt"

type asciiCode string

func (a asciiCode) str() string {
	return string(a)
}

const (
	reset       asciiCode = "\033[0m"
	resetColor  asciiCode = "\033[39m"
	symbolColor asciiCode = "\033[90m"

	bold        asciiCode = "\033[1m"
	italic      asciiCode = "\033[3m"
	boldItalic  asciiCode = "\033[1m\033[3m"
	underline   asciiCode = "\033[4m"
	stricked    asciiCode = "\033[9m"
	quote       asciiCode = "\033[32m"
	quoteSymbol asciiCode = "\u2503"

	startSel asciiCode = "\033[100m"
	endSel   asciiCode = "\033[49m"

	codeLine asciiCode = "\033[48;5;236m\033[33m"
	header   asciiCode = "\033[34m"

	listColor asciiCode = "\033[35m"
	tagColor  asciiCode = "\033[48;5;60m\033[38;5;219m"
	tagS                = "["
	tagE                = "]"
	shield    asciiCode = "\\"
	listDash  asciiCode = "\u2981 "
	boxEmpty  asciiCode = " \u25a1"
	boxField  asciiCode = " \u25a0"
)

type Renderer struct {
	curAttr asciiCode
	w, h    int
	l       *lexer
}

func InitReder(w, h int) *Renderer {
	r := &Renderer{
		w: w,
		h: h,
	}
	//TODO: create a new lexer, for code, and separate it form default markdown lexer and renderer
	r.l = NewLexer()
	return r
}

func (r *Renderer) RednerMarkdownLine(line []rune, isCur bool) (string, int) {
	//here is reset lexer, so it can read a new line. Prev I was creating a new instance of a lexer, which is not rly good, ig

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
			data += r.simpleAttrRender(codeLine, string(tok.Literal), isCur)
			diff += 1
		case TEXT:
			data += r.renderText(&tok)
		case Shield:
			data += r.renderShield(&tok, isCur)
			diff += 1
		case Tag:
			data += r.renderTag(&tok, isCur)
			diff -= 1
		case OneStar, OneUnderline:
			data += r.simpleAttrRender(italic, string(tok.Literal), isCur)
			diff += 1
		case TwoStars, TwoUnderlines:
			data += r.simpleAttrRender(bold, string(tok.Literal), isCur)
			diff += 2
		case ThreeStars, ThreeUnderlines:
			data += r.simpleAttrRender(boldItalic, string(tok.Literal), isCur)
			diff += 3
		case Stricked:
			data += r.simpleAttrRender(stricked, string(tok.Literal), isCur)
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
	return sym + resetColor.str()
}

func paintString(ascii asciiCode, str string) string {
	var s = ""
	for _, x := range str {
		s += fmt.Sprintf("%s%c", ascii, x)
	}
	return s
}

func (r *Renderer) renderBoxEmpty(t *Token, isCur bool) string {
	var s = ""
	if isCur {
		s += painAsAttr(string(t.Literal))
	} else {
		s += boxEmpty.str()
	}
	return s
}

func (r *Renderer) renderBoxField(t *Token, isCur bool) string {
	var s = ""
	if isCur {
		s += painAsAttr(string(t.Literal))
	} else {
		s += boxField.str()
	}
	return s
}

func (r *Renderer) renderListNumber(t *Token) string {
	var s = ""
	s += paintString(listColor, string(t.Value))
	s += paintString(listColor, string(t.Literal))
	s += reset.str()
	return s
}

func (r *Renderer) renderListDash(t *Token, isCur bool) string {
	var s = ""
	if isCur {
		s += painAsAttr(string(t.Literal))
	} else {
		s += listDash.str()
	}
	return s
}

func (r *Renderer) renderShield(t *Token, isCur bool) string {
	var s = ""
	if isCur {
		s += painAsAttr(string(t.Literal))
	}
	s += string(t.Value)
	return s
}

func (r *Renderer) renderQuote(t *Token, isCur bool) string {
	var s = ""
	s += quote.str()
	if isCur {
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

func (r *Renderer) renderTag(t *Token, isCur bool) string {
	var s = ""
	if !isCur {
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

func (r *Renderer) simpleAttrRender(mode asciiCode, attr string, isCur bool) string {
	var s = ""
	if r.curAttr == mode {
		r.curAttr = reset
		if isCur {
			s += painAsAttr(attr)
		}
		s += r.curAttr.str()
	} else {
		r.curAttr = mode
		s += r.curAttr.str()
		if isCur {
			s += painAsAttr(attr)
		}
	}
	return s
}
