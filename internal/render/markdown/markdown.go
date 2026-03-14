package markdown

import (
	"strings"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/render/general"
)

type Render struct {
	curAttr string
	w       int
	l       *Lexer
}

func NewRender(w int) *Render {
	r := &Render{w: w}
	r.l = newLexer()
	return r
}

func (r *Render) Reset() {
	r.l.input = []rune{}
	r.l.position = 0
	r.l.readPosition = 0
}

func (r *Render) RenderMarkdownLine(line []rune, lineIndex int, show bool) (string, int, enums.Render) {
	var renderMode enums.Render = enums.Markdown

	if string(line) == "---" || string(line) == "***" || string(line) == "___" {
		if show {
			return general.PainAsAttr("---"), 0, renderMode
		}
		return general.PainAsAttr(strings.Repeat("\u2015", r.w)), 3 - r.w + enums.InitialOffset*2, renderMode
	}

	r.l.input = line
	r.l.position = 0
	r.l.readPosition = 0
	r.l.readChar()

	isWhiteSpace := true

	var data strings.Builder
	var diff = 0

	var i = 0
	for tok := r.l.NextToken(); tok.Type != EOL; tok = r.l.NextToken() {
		switch tok.Type {
		case Header_1, Header_2, Header_3, Header_4, Header_5, Header_6:
			if i == 0 {
				data.WriteString(r.renderHeader(&tok))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case ListBoxField:
			if i == 0 {
				data.WriteString(r.renderBoxField(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case ListBoxEmpty:
			if i == 0 {
				data.WriteString(r.renderBoxEmpty(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case Quote:
			if i == 0 {
				data.WriteString(r.renderQuote(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case CodeBlock:
			if i == 0 {
				data.WriteString(r.renderCodeBlock(&tok, show))
			} else {
				data.WriteString(string(tok.Literal) + string(tok.Value))
			}
			renderMode = enums.Code
		case ListDash:
			if isWhiteSpace {
				data.WriteString(r.renderListDash(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case ListNumberB, ListNumberDot:
			if isWhiteSpace {
				data.WriteString(r.renderListNumber(&tok))
			} else {
				data.WriteString(string(tok.Value) + string(tok.Literal))
			}
		case Hightlight:
			data.WriteString(r.simpleAttrRender(ascii.Hightlight.Str(), string(tok.Value), show))
		case Link:
			data.WriteString(r.renderLink(&tok, show))
		case Image:
			data.WriteString(r.renderImage(&tok, show))
		case CodeLine:
			data.WriteString(r.simpleAttrRender(ascii.CodeLine.Str(), string(tok.Literal), show))
		case TEXT:
			data.WriteString(r.renderText(&tok))
		case Shield:
			data.WriteString(r.renderShield(&tok, show))
		case Tag:
			data.WriteString(r.renderTag(&tok, show))
			diff -= 1
		case OneStar, OneUnderline:
			data.WriteString(r.simpleAttrRender(ascii.Italic.Str(), string(tok.Literal), show))
		case TwoStars, TwoUnderlines:
			data.WriteString(r.simpleAttrRender(ascii.Bold.Str(), string(tok.Literal), show))
		case ThreeStars, ThreeUnderlines:
			data.WriteString(r.simpleAttrRender(ascii.BoldItalic.Str(), string(tok.Literal), show))
		case Stricked:
			data.WriteString(r.simpleAttrRender(ascii.Stricked.Str(), string(tok.Literal), show))
		case WhiteSpace:
			data.WriteString(" ")
		case Symbol:
			data.WriteString(string(tok.Value))
		}
		i += 1
		if isWhiteSpace && tok.Type != WhiteSpace {
			isWhiteSpace = false
		}
	}

	data.WriteString(ascii.Reset.Str())
	r.curAttr = ascii.Reset.Str()
	// if !show {
	// 	r.c.cacheLine(line, data, diff, lineIndex)
	// }
	return data.String(), 0, renderMode
}

func (r *Render) renderCodeBlock(t *Token, show bool) string {
	//r.change = true
	if show {
		return general.PainAsAttr(string(t.Literal)) + general.PaintString(ascii.CodeLine, string(t.Value))
	} else {
		return general.PaintString(ascii.CodeLine, string(t.Value))
	}
}

func (r *Render) renderBoxEmpty(t *Token, show bool) string {
	if show {
		return general.PainAsAttr(string(t.Literal))
	} else {
		return ascii.BoxEmpty.Str()
	}
}

func (r *Render) renderBoxField(t *Token, show bool) string {
	if show {
		return general.PainAsAttr(string(t.Literal))
	} else {
		return ascii.BoxField.Str()
	}
}

func (r *Render) renderListNumber(t *Token) string {
	var s strings.Builder
	s.WriteString(general.PaintString(ascii.ListColor, string(t.Value)))
	s.WriteString(general.PaintString(ascii.ListColor, string(t.Literal)))
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderListDash(t *Token, show bool) string {
	if show {
		return general.PainAsAttr(string(t.Literal))
	} else {
		return ascii.ListDash.Str()
	}
}

func (r *Render) renderShield(t *Token, show bool) string {
	var s strings.Builder
	if show {
		s.WriteString(general.PainAsAttr(string(t.Literal)))
	}
	s.WriteString(string(t.Value))
	return s.String()
}

func (r *Render) renderQuote(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(ascii.Quote.Str())
	if show {
		s.WriteString(general.PainAsAttr(string(t.Literal)))
	} else {
		s.WriteString(ascii.QuoteSymbol.Str())
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderText(t *Token) string {
	if r.curAttr != ascii.Reset.Str() {
		return general.PaintString(ascii.Color(r.curAttr), string(t.Value))
	}
	return string(t.Value)
}

func (r *Render) renderTag(t *Token, show bool) string {
	var s strings.Builder
	if !show {
		s.WriteString(ascii.TagColor.Str())
		s.WriteString(ascii.TagS.Str())
		s.WriteString(general.PaintString(ascii.TagColor, string(t.Literal)))
		s.WriteString(general.PaintString(ascii.TagColor, string(t.Value)))
		s.WriteString(ascii.TagE.Str())
		s.WriteString(ascii.Reset.Str())
	} else {
		s.WriteString(general.PaintString(ascii.TagColor, string(t.Literal)))
		s.WriteString(general.PaintString(ascii.TagColor, string(t.Value)))
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderHeader(t *Token) string {
	var s strings.Builder
	s.WriteString(ascii.Header.Str())
	s.WriteString(ascii.Underline.Str())
	r.curAttr = ascii.Header.Str()
	s.WriteString(string(t.Literal))
	return s.String()
}

func (r *Render) renderLink(t *Token, show bool) string {
	if show {
		return ascii.Link.Str() + string(t.Literal) + ascii.Reset.Str()
	} else {
		return ascii.Link.Str() + string(t.Value) + ascii.Reset.Str()
	}
}

func (r *Render) renderImage(t *Token, show bool) string {
	if show {
		return ascii.Link.Str() + string(t.Literal) + ascii.Reset.Str()
	} else {
		return ascii.Link.Str() + string(t.Value) + ascii.Reset.Str()
	}
}

func (r *Render) simpleAttrRender(mode string, attr string, show bool) string {
	var s strings.Builder
	if r.curAttr == mode {
		r.curAttr = ascii.Reset.Str()
		if show {
			s.WriteString(general.PainAsAttr(attr))
		}
		s.WriteString(r.curAttr)
	} else {
		r.curAttr = mode
		s.WriteString(r.curAttr)
		if show {
			s.WriteString(general.PainAsAttr(attr))
		}
	}
	return s.String()
}
