package markdown

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"
	"github.com/EnotInc/Bard/theme"

	render "github.com/EnotInc/Bard/internal/enums/render"
)

type Render struct {
	l       *Lexer
	theme   *theme.Markdown
	curAttr string
	w       int
}

func (r *Render) Update() {
	theme := theme.GetTheme().Markdown
	r.theme = &theme
}

func NewRender(w int) *Render {
	theme := theme.GetTheme().Markdown
	r := &Render{w: w, theme: &theme}
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

func isLine(line []rune, key rune) bool {
	return len(line) == 3 && line[0] == line[1] && line[1] == line[2] && line[2] == key
}

func (r *Render) RenderMarkdownLine(line []rune, lineIndex int, show bool, xOffset int) (string, render.Render, bool) {
	var renderMode render.Render = render.Markdown

	if isLine(line, '-') || isLine(line, '*') || isLine(line, '_') {
		var data strings.Builder
		var keep bool

		if show {
			data.WriteString(r.theme.Symbol)
			data.WriteString(string(line))
			keep = false
		} else {
			data.WriteString(r.theme.Symbol)
			data.WriteString(strings.Repeat(ascii.SplitLine.Str(), r.w-enums.InitialOffset*2+xOffset))
			data.WriteString(ascii.Reset.Str())
			keep = true
		}
		return data.String(), renderMode, keep
	}

	r.l.input = line
	r.l.position = 0
	r.l.readPosition = 0
	r.l.readChar()

	var data strings.Builder

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
		case listBoxComplete:
			if isWhiteSpace {
				data.WriteString(r.renderBoxComplete(&tok, show))
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
		case quote_note:
			if isWhiteSpace {
				data.WriteString(r.renderQuoteNote(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case quote_tip:
			if isWhiteSpace {
				data.WriteString(r.renderQuoteTip(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case quote_important:
			if isWhiteSpace {
				data.WriteString(r.renderQuoteImportant(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case quote_warning:
			if isWhiteSpace {
				data.WriteString(r.renderQuoteWarning(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case quote_caution:
			if isWhiteSpace {
				data.WriteString(r.renderQuoteCaution(&tok, show))
			} else {
				data.WriteString(string(tok.Literal))
			}
		case codeBlock:
			if isFirst {
				data.WriteString(r.renderCodeBlock(&tok, show, xOffset))
				renderMode = render.Code
			} else {
				data.WriteString(string(tok.Literal))
				data.WriteString(string(tok.Value))
			}
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
				data.WriteString(string(tok.Value))
				data.WriteString(string(tok.Literal))
			}
		case tab:
			data.WriteString(r.renderTab(&tok))
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
		case colorBlock:
			data.WriteString(r.renderColorBlock(&tok, show))
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
		if isWhiteSpace && (tok.Type != whitespace && tok.Type != tab) {
			isWhiteSpace = false
		}
	}

	data.WriteString(ascii.Reset.Str())
	r.curAttr = ascii.Reset.Str()
	return data.String(), renderMode, false
}

func (r *Render) fillSpace(xScroll int) string {
	amount := max(r.w-len(r.l.input)-enums.InitialOffset, 0)
	return strings.Repeat(" ", amount+xScroll)
}

func (r *Render) renderWSEOL(t *Token) string {
	return strings.Repeat(r.theme.Symbol+ascii.WSEOL.Str(), len(t.Value))
}

func (r *Render) renderCodeBlock(t *Token, show bool, xScroll int) string {
	if show {
		return r.theme.CodeHeader + r.theme.Symbol + string(t.Literal) + string(t.Value) + r.fillSpace(xScroll)
	}

	i := ""
	si := config.GetConfig().ShowIcons
	if len(t.Value) > 0 {
		i = services.GetFileIcon(string(t.Value), si)
	}

	return r.theme.CodeHeader + " " + i + r.theme.Symbol + string(t.Value) + r.fillSpace(xScroll)
}

func (r *Render) renderBoxEmpty(t *Token, show bool) string {
	if show {
		return r.theme.Symbol + string(t.Literal) + ascii.Reset.Str()
	}
	return ascii.BoxEmpty.Str()
}

func (r *Render) renderBoxComplete(t *Token, show bool) string {
	if show {
		return r.theme.Symbol + string(t.Literal) + ascii.Reset.Str()
	}
	return ascii.BoxComplete.Str()
}

func (r *Render) renderBoxFilled(t *Token, show bool) string {
	if show {
		return r.theme.Symbol + string(t.Literal) + ascii.Reset.Str()
	}
	return ascii.BoxField.Str()
}

func (r *Render) renderListNumber(t *Token) string {
	var s strings.Builder
	s.WriteString(r.theme.NumberList)
	s.WriteString(string(t.Value))
	s.WriteString(string(t.Literal))
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderListDash(t *Token, show bool) string {
	if show {
		return r.theme.Symbol + string(t.Literal) + ascii.Reset.Str()
	}
	return ascii.ListDash.Str()
}

func (r *Render) renderShield(t *Token, show bool) string {
	var s strings.Builder
	if show {
		s.WriteString(r.theme.Symbol)
		s.WriteString(string(t.Literal))
	}
	s.WriteString(string(t.Value))
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderQuote(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(r.theme.Quote)
	if show {
		s.WriteString(r.theme.Symbol)
		s.WriteString(string(t.Literal))
	} else {
		quote := strings.Repeat(ascii.QuoteSymbol.Str(), len(t.Literal))
		s.WriteString(quote)
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderQuoteNote(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(r.theme.Quote)
	if show {
		s.WriteString(r.theme.Symbol)
		s.WriteString(string(t.Literal))
	} else {
		spl := strings.Split(string(t.Literal), " ")
		quote := strings.Repeat(ascii.QuoteSymbol.Str(), len(spl[0]))
		s.WriteString(quote)
		s.WriteString(" ")
		s.WriteString(ascii.QuoteNote.Str())
		s.WriteString(" Note")
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderQuoteTip(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(r.theme.Quote)
	if show {
		s.WriteString(r.theme.Symbol)
		s.WriteString(string(t.Literal))
	} else {
		spl := strings.Split(string(t.Literal), " ")
		quote := strings.Repeat(ascii.QuoteSymbol.Str(), len(spl[0]))
		s.WriteString(quote)
		s.WriteString(ascii.QuoteTip.Str())
		s.WriteString(" Tip")
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderQuoteImportant(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(r.theme.Quote)
	if show {
		s.WriteString(r.theme.Symbol)
		s.WriteString(string(t.Literal))
	} else {
		spl := strings.Split(string(t.Literal), " ")
		quote := strings.Repeat(ascii.QuoteSymbol.Str(), len(spl[0]))
		s.WriteString(quote)
		s.WriteString(" ")
		s.WriteString(ascii.QuoteImportant.Str())
		s.WriteString(" Important")
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderQuoteWarning(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(r.theme.Quote)
	if show {
		s.WriteString(r.theme.Symbol)
		s.WriteString(string(t.Literal))
	} else {
		spl := strings.Split(string(t.Literal), " ")
		quote := strings.Repeat(ascii.QuoteSymbol.Str(), len(spl[0]))
		s.WriteString(quote)
		s.WriteString(" ")
		s.WriteString(ascii.QuoteWarning.Str())
		s.WriteString(" Warning")
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderQuoteCaution(t *Token, show bool) string {
	var s strings.Builder
	s.WriteString(r.theme.Quote)
	if show {
		s.WriteString(r.theme.Symbol)
		s.WriteString(string(t.Literal))
	} else {
		spl := strings.Split(string(t.Literal), " ")
		quote := strings.Repeat(ascii.QuoteSymbol.Str(), len(spl[0]))
		s.WriteString(quote)
		s.WriteString(ascii.QuoteCaution.Str())
		s.WriteString(" Caution")
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderText(t *Token) string {
	if r.curAttr != ascii.Reset.Str() {
		return r.curAttr + string(t.Value)
	}
	return string(t.Value)
}

func (r *Render) renderTag(t *Token, show bool) string {
	var s strings.Builder
	if !show {
		s.WriteString(r.theme.Tag)
		s.WriteString(ascii.TagS.Str())
		s.WriteString(r.theme.Tag)
		s.WriteString(string(t.Literal))
		s.WriteString(string(t.Value))
		s.WriteString(ascii.TagE.Str())
		s.WriteString(ascii.Reset.Str())
	} else {
		s.WriteString(r.theme.Tag)
		s.WriteString(string(t.Literal))
		s.WriteString(string(t.Value))
	}
	s.WriteString(ascii.Reset.Str())
	return s.String()
}

func (r *Render) renderColorBlock(t *Token, show bool) string {
	var s strings.Builder

	color, err := services.HexToAscii(string(t.Value), services.Foreground)
	if err != nil {
		s.WriteString(string(t.Literal))
		s.WriteString(string(t.Value))
		return s.String()
	}

	s.WriteString(color)

	if !show {
		s.WriteString(ascii.ColorBox.Str())
	} else {
		s.WriteString(string(t.Literal))
	}

	s.WriteString(ascii.ResetFg.Str())
	s.WriteString(string(t.Value))
	return s.String()
}

func (r *Render) renderHeader(t *Token) string {
	var header strings.Builder
	switch t.Type {
	case header_1:
		header.WriteString(r.theme.Header1)
	case header_2:
		header.WriteString(r.theme.Header2)
	case header_3:
		header.WriteString(r.theme.Header3)
	case header_4:
		header.WriteString(r.theme.Header4)
	case header_5:
		header.WriteString(r.theme.Header5)
	case header_6:
		header.WriteString(r.theme.Header6)
	}
	header.WriteString(string(t.Literal))
	return header.String()
}

func (r *Render) renderTab(t *Token) string {
	if len(t.Literal) > 0 {
		return r.theme.Symbol + ascii.Tab.Str() + ascii.ResetFg.Str() + string(t.Literal[1:])
	} else {
		return r.theme.Symbol + ascii.Tab.Str() + ascii.ResetFg.Str()
	}
}

func (r *Render) renderLink(t *Token, show bool) string {
	var data strings.Builder
	if show {
		data.WriteString(r.theme.Link)
		data.WriteString(string(t.Literal))
	} else {
		data.WriteString(r.theme.Link)
		data.WriteString(ascii.UnderLine.Str())
		data.WriteString(ascii.LinkSymbol.Str())
		data.WriteString(" ")
		data.WriteString(string(t.Value))
	}
	data.WriteString(ascii.Reset.Str())
	return data.String()
}

func (r *Render) renderImage(t *Token, show bool) string {
	var data strings.Builder
	if show {
		data.WriteString(r.theme.Image)
		data.WriteString(string(t.Literal))
	} else {
		data.WriteString(r.theme.Image)
		data.WriteString(ascii.UnderLine.Str())
		data.WriteString(ascii.ImageSymbol.Str())
		data.WriteString(" ")
		data.WriteString(string(t.Value))
	}
	data.WriteString(ascii.Reset.Str())
	return data.String()
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
		s.WriteString(string(t.Literal))
		s.WriteString(string(t.Value))
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
			s.WriteString(r.theme.Symbol)
			s.WriteString(attr)
		}
		s.WriteString(r.curAttr)
	} else {
		r.curAttr = mode
		s.WriteString(r.curAttr)
		if show {
			s.WriteString(r.theme.Symbol)
			s.WriteString(attr)
		}
	}
	return s.String()
}
