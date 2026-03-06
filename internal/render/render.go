package render

import (
	"Enot/Bard/internal/ascii"
	"Enot/Bard/internal/enums"
	"fmt"
	"slices"
	"strings"
)

type Renderer struct {
	curAttr string
	c       *cache
	l       *lexer
	w       int
}

func InitRender(w, h int) *Renderer {
	_c := initCache()
	r := &Renderer{c: _c, w: w}
	// TODO: create a new lexer for code and separate it from the default markdown lexer and renderer
	r.l = NewLexer()
	return r
}

func (r *Renderer) RenderMarkdownLine(line []rune, lineIndex int, show bool) (string, int) {

	if !show && r.c.isCached(lineIndex) {
		l := r.c.getCached(lineIndex)
		if slices.Equal(l.raw, line) {
			return l.render, l.diff
		}
	}

	if string(line) == "---" || string(line) == "***" || string(line) == "___" {
		if show {
			return painAsAttr("---"), 0
		}
		return painAsAttr(strings.Repeat("\u2015", r.w)), 3 - r.w + enums.InitialOffset*2
	}

	r.l.input = line
	r.l.position = 0
	r.l.readPosition = 0
	r.l.readChar()

	isWhiteSpace := true

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
		case Quote:
			if i == 0 {
				data += r.renderQuote(&tok, show)
			} else {
				data += string(tok.Literal)
			}
		case ListDash:
			if isWhiteSpace {
				data += r.renderListDash(&tok, show)
			} else {
				data += string(tok.Literal)
			}
		case ListNumberB, ListNumberDot:
			if isWhiteSpace {
				data += r.renderListNumber(&tok)
			} else {
				data += string(tok.Value) + string(tok.Literal)
			}
		case Hightlight:
			data += r.simpleAttrRender(ascii.Hightlight.Str(), string(tok.Value), show)
		case Link:
			data += r.renderLink(&tok, show)
		case Image:
			data += r.renderImage(&tok, show)
		case CodeLine:
			data += r.simpleAttrRender(ascii.CodeLine.Str(), string(tok.Literal), show)
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
			data += r.simpleAttrRender(ascii.Italic.Str(), string(tok.Literal), show)
			diff += 1
		case TwoStars, TwoUnderlines:
			data += r.simpleAttrRender(ascii.Bold.Str(), string(tok.Literal), show)
			diff += 2
		case ThreeStars, ThreeUnderlines:
			data += r.simpleAttrRender(ascii.BoldItalic.Str(), string(tok.Literal), show)
			diff += 3
		case Stricked:
			data += r.simpleAttrRender(ascii.Stricked.Str(), string(tok.Literal), show)
			diff += 2
		case WhiteSpace:
			data += " "
		case Symbol:
			data += string(tok.Value)
		}
		i += 1
		if isWhiteSpace && tok.Type != WhiteSpace {
			isWhiteSpace = false
		}
	}

	data += ascii.Reset.Str()
	r.curAttr = ascii.Reset.Str()
	if !show {
		r.c.cacheLine(line, data, diff, lineIndex)
	}
	return data, diff
}

func painAsAttr(symbol string) string {
	sym := paintString(ascii.SymbolColor, symbol)
	return sym + ascii.Reset.Str()
}

func paintString(c ascii.Color, str string) string {
	var s = ""
	for _, x := range str {
		s += fmt.Sprintf("%s%c", c, x)
	}
	return s
}

func (r *Renderer) renderBoxEmpty(t *Token, show bool) string {
	if show {
		return painAsAttr(string(t.Literal))
	} else {
		return ascii.BoxEmpty.Str()
	}
}

func (r *Renderer) renderBoxField(t *Token, show bool) string {
	if show {
		return painAsAttr(string(t.Literal))
	} else {
		return ascii.BoxField.Str()
	}
}

func (r *Renderer) renderListNumber(t *Token) string {
	var s = ""
	s += paintString(ascii.ListColor, string(t.Value))
	s += paintString(ascii.ListColor, string(t.Literal))
	s += ascii.Reset.Str()
	return s
}

func (r *Renderer) renderListDash(t *Token, show bool) string {
	if show {
		return painAsAttr(string(t.Literal))
	} else {
		return ascii.ListDash.Str()
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
	s += ascii.Quote.Str()
	if show {
		s += painAsAttr(string(t.Literal))
	} else {
		s += ascii.QuoteSymbol.Str()
	}
	s += ascii.Reset.Str()
	return s
}

func (r *Renderer) renderText(t *Token) string {
	if r.curAttr != ascii.Reset.Str() {
		return paintString(ascii.Color(r.curAttr), string(t.Value))
	}
	return string(t.Value)
}

func (r *Renderer) renderTag(t *Token, show bool) string {
	var s = ""
	if !show {
		s += ascii.TagColor.Str()
		s += ascii.TagS.Str()
		s += paintString(ascii.TagColor, string(t.Literal))
		s += paintString(ascii.TagColor, string(t.Value))
		s += ascii.TagE.Str()
		s += ascii.Reset.Str()
	} else {
		s += paintString(ascii.TagColor, string(t.Literal))
		s += paintString(ascii.TagColor, string(t.Value))
	}
	s += ascii.Reset.Str()
	return s
}

func (r *Renderer) renderHeader(t *Token) string {
	var s = ""
	s += ascii.Header.Str()
	s += ascii.Underline.Str()
	r.curAttr = ascii.Header.Str()
	s += string(t.Literal)
	return s
}

func (r *Renderer) renderLink(t *Token, show bool) string {
	if show {
		return ascii.Link.Str() + string(t.Literal) + ascii.Reset.Str()
	} else {
		return ascii.Link.Str() + string(t.Value) + ascii.Reset.Str()
	}
}

func (r *Renderer) renderImage(t *Token, show bool) string {
	if show {
		return ascii.Link.Str() + string(t.Literal) + ascii.Reset.Str()
	} else {
		return ascii.Link.Str() + string(t.Value) + ascii.Reset.Str()
	}
}

func (r *Renderer) simpleAttrRender(mode string, attr string, show bool) string {
	var s = ""
	if r.curAttr == mode {
		r.curAttr = ascii.Reset.Str()
		if show {
			s += painAsAttr(attr)
		}
		s += r.curAttr
	} else {
		r.curAttr = mode
		s += r.curAttr
		if show {
			s += painAsAttr(attr)
		}
	}
	return s
}
