package markdown

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/render/general"
)

type Render struct {
	curAttr string
	w       int
	l       *Lexer
	theme   *config.Markdown
}

func NewRender(w int, theme *config.Markdown) *Render {
	r := &Render{w: w, theme: theme}
	r.l = newLexer()
	return r
}

func (r *Render) Resize(w int) {
	r.w = w
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
			return general.PaintString(r.theme.Symbol, string(line)), 0, renderMode
		}
		return general.PaintString(r.theme.Symbol, strings.Repeat(ascii.SplitLIne.Str(), r.w)), 3 - r.w + enums.InitialOffset*2, renderMode
	}

	r.l.input = line
	r.l.position = 0
	r.l.readPosition = 0
	r.l.readChar()

	var data strings.Builder
	var diff = 0

	isWhiteSpace := true
	isFirst := true

	for tok := r.l.NextToken(); tok.Type != eol; tok = r.l.NextToken() {
		switch tok.Type {
		case header_1, header_2, header_3, header_4, header_5, header_6:
			if isFirst {
				data.WriteString(r.renderHeader(&tok))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case listBoxFilled:
			if isWhiteSpace {
				data.WriteString(r.renderBoxFilled(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case listBoxEmpty:
			if isWhiteSpace {
				data.WriteString(r.renderBoxEmpty(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case quote:
			if isWhiteSpace {
				data.WriteString(r.renderQuote(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case codeBlock:
			if isWhiteSpace {
				data.WriteString(r.renderCodeBlock(&tok, show))
				diff = -r.w - len(r.l.input)
			} else {
				data.WriteString(string(tok.Literal) + string(tok.Value))
			}
			renderMode = enums.Code
		case listDash:
			if isWhiteSpace {
				data.WriteString(r.renderListDash(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case listNumberB, listNumberDot:
			if isWhiteSpace {
				data.WriteString(r.renderListNumber(&tok))
			} else {
				data.WriteString(string(tok.Value) + string(tok.Literal))
			}
		case tab:
			data.WriteString(r.renderTab(&tok))
			diff -= len(tok.Literal)
		case hightlight:
			data.WriteString(r.simpleAttrRender(r.theme.Highlight, string(tok.Value), show))
		case link:
			data.WriteString(r.renderLink(&tok, show))
		case image:
			data.WriteString(r.renderImage(&tok, show))
		case html:
			data.WriteString(r.renderHtmlBlcok(&tok))
		case codeLine:
			data.WriteString(r.renderCodeLine(&tok, show))
		case text:
			data.WriteString(r.renderText(&tok))
		case shield:
			data.WriteString(r.renderShield(&tok, show))
		case tag:
			data.WriteString(r.renderTag(&tok, show))
			diff -= 1
		case oneStar, oneUnderLine:
			data.WriteString(r.simpleAttrRender(ascii.Italic.Str(), string(tok.Literal), show))
		case twoStars, twoUnderLines:
			data.WriteString(r.simpleAttrRender(ascii.Bold.Str(), string(tok.Literal), show))
		case threeStars, threeUnderLines:
			data.WriteString(r.simpleAttrRender(ascii.BoldItalic.Str(), string(tok.Literal), show))
		case stricked:
			data.WriteString(r.simpleAttrRender(ascii.Stricked.Str(), string(tok.Literal), show))
		case whitespace:
			data.WriteString(string(tok.Value))
		case wseol:
			data.WriteString(r.renderWSEOL(&tok))
		case symbol:
			data.WriteString(string(tok.Value))
		}
		isFirst = false
		if isWhiteSpace && tok.Type != whitespace {
			isWhiteSpace = false
		}
	}

	data.WriteString(ascii.Reset.Str())
	r.curAttr = ascii.Reset.Str()
	return data.String(), diff, renderMode
}

func (r *Render) fillSpace() string {
	amount := max(r.w-len(r.l.input)-enums.InitialOffset-1, 0)
	return strings.Repeat(" ", amount)
}

func (r *Render) renderWSEOL(t *Token) string {
	return strings.Repeat(r.theme.Symbol+ascii.WSEOL.Str(), len(t.Value))
}

func (r *Render) renderCodeBlock(t *Token, show bool) string {
	if show {
		return r.theme.CodeHeader + general.PaintString(r.theme.Symbol, string(t.Literal)+string(t.Value)) + r.fillSpace()
	}

	if i, ok := langIcon[strings.ToLower(string(t.Value))]; ok {
		return r.theme.CodeHeader + " " + i + general.PaintString(r.theme.Symbol, string(t.Value)) + r.fillSpace()
	}
	// fallback
	return r.theme.CodeHeader + general.PaintString(r.theme.Symbol, "["+string(t.Value)+"] ") + r.fillSpace()
}

func (r *Render) renderBoxEmpty(t *Token, show bool) string {
	if show {
		return general.PaintString(r.theme.Symbol, string(t.Literal)) + ascii.Reset.Str()
	}
	return ascii.BoxEmpty.Str()
}

func (r *Render) renderBoxFilled(t *Token, show bool) string {
	if show {
		return general.PaintString(r.theme.Symbol, string(t.Literal)) + ascii.Reset.Str()
	}
	return ascii.BoxField.Str()
}

func (r *Render) renderListNumber(t *Token) string {
	var s strings.Builder
	s.WriteString(general.PaintString(r.theme.NumberList, string(t.Value)))
	s.WriteString(general.PaintString(r.theme.NumberList, string(t.Literal)))
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderListDash(t *Token, show bool) string {
	if show {
		return general.PaintString(r.theme.Symbol, string(t.Literal)) + ascii.Reset.Str()
	}
	return ascii.ListDash.Str()
}

func (r *Render) renderShield(t *Token, show bool) string {
	var s strings.Builder
	if show {
		s.WriteString(general.PaintString(r.theme.Symbol, string(t.Literal)))
	}
	s.WriteString(string(t.Value))
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderQuote(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(r.theme.Quote)
	if show {
		s.WriteString(general.PaintString(r.theme.Symbol, string(t.Literal)))
	} else {
		s.WriteString(ascii.QuoteSymbol.Str())
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderText(t *Token) string {
	if r.curAttr != ascii.Reset.Str() {
		return general.PaintString(r.curAttr, string(t.Value))
	}
	return string(t.Value)
}

func (r *Render) renderTag(t *Token, show bool) string {
	var s strings.Builder
	if !show {
		s.WriteString(r.theme.Tag)
		s.WriteString(ascii.TagS.Str())
		s.WriteString(general.PaintString(r.theme.Tag, string(t.Literal)))
		s.WriteString(general.PaintString(r.theme.Tag, string(t.Value)))
		s.WriteString(ascii.TagE.Str())
		s.WriteString(ascii.Reset.Str())
	} else {
		s.WriteString(general.PaintString(r.theme.Tag, string(t.Literal)))
		s.WriteString(general.PaintString(r.theme.Tag, string(t.Value)))
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderHeader(t *Token) string {
	var s strings.Builder
	header := r.theme.Header1
	switch t.Type {
	case header_1:
		header = r.theme.Header1
	case header_2:
		header = r.theme.Header2
	case header_3:
		header = r.theme.Header3
	case header_4:
		header = r.theme.Header4
	case header_5:
		header = r.theme.Header5
	case header_6:
		header = r.theme.Header6
	}
	s.WriteString(header)
	r.curAttr = header
	s.WriteString(string(t.Literal))
	return s.String()
}

func (r *Render) renderTab(t *Token) string {
	return r.theme.Symbol + ascii.Tab.Str() + ascii.Reset.Str() + string(t.Literal[1:])
}

func (r *Render) renderLink(t *Token, show bool) string {
	if show {
		return r.theme.Link + string(t.Literal) + ascii.Reset.Str()
	}
	return r.theme.Link + string(t.Value) + ascii.Reset.Str()
}

func (r *Render) renderImage(t *Token, show bool) string {
	if show {
		return r.theme.Image + string(t.Literal) + ascii.Reset.Str()
	}
	return r.theme.Image + string(t.Value) + ascii.Reset.Str()
}

func (r *Render) renderHtmlBlcok(t *Token) string {
	var s strings.Builder
	switch len(t.Literal) {
	case 3: // </>
		s.WriteString(r.theme.HTMLSymbol)
		s.WriteString(string(t.Literal[:len(t.Literal)-1]))
		s.WriteString(r.theme.HTMLText)
		s.WriteString(string(t.Value))
		s.WriteString(r.theme.HTMLSymbol)
		s.WriteString(string(t.Literal[len(t.Literal)-1]))
	case 2: // </ or <>
		if t.Literal[1] == '/' {
			s.WriteString(r.theme.HTMLSymbol)
			s.WriteString(string(t.Literal))
			s.WriteString(r.theme.HTMLText)
			s.WriteString(string(t.Value))
		} else {
			s.WriteString(r.theme.HTMLSymbol)
			s.WriteString(string(t.Literal[:len(t.Literal)-1]))
			s.WriteString(r.theme.HTMLText)
			s.WriteString(string(t.Value))
			s.WriteString(r.theme.HTMLSymbol)
			s.WriteString(string(t.Literal[len(t.Literal)-1]))
		}
	default: // <
		s.WriteString(r.theme.HTMLSymbol)
		s.WriteString(string(t.Literal))
		s.WriteString(r.theme.HTMLText)
		s.WriteString(string(t.Value))
	}
	return s.String() + ascii.Reset.Str()
}

func (r *Render) renderCodeLine(t *Token, show bool) string {
	var s strings.Builder
	if show {
		s.WriteString(string(t.Literal) + string(t.Value))
	} else {
		end := len(t.Value)
		if end > 0 {
			end -= 1
		}
		s.WriteString(string(t.Value[:end]))
	}
	return r.theme.CodeLineBg + r.theme.CodeText + s.String() + ascii.Reset.Str()
}

func (r *Render) simpleAttrRender(mode string, attr string, show bool) string {
	var s strings.Builder
	if r.curAttr == mode {
		r.curAttr = ascii.Reset.Str()
		if show {
			s.WriteString(general.PaintString(r.theme.Symbol, attr))
		}
		s.WriteString(r.curAttr)
	} else {
		r.curAttr = mode
		s.WriteString(r.curAttr)
		if show {
			s.WriteString(general.PaintString(r.theme.Symbol, attr))
		}
	}
	return s.String()
}
