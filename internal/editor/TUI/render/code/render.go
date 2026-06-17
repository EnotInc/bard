package code

import (
	"strings"

	"github.com/EnotInc/Bard/config"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/enums/ascii"
	"github.com/EnotInc/Bard/internal/services"

	render "github.com/EnotInc/Bard/internal/enums/render"
)

type Render struct {
	l     *Lexer
	theme *config.Code
	w     int
}

func (r *Render) Update() {
	theme := config.GetTheme().Code
	r.theme = &theme
}

func NewRender(w int) *Render {
	theme := config.GetTheme().Code
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

func (r *Render) fillSpace(xScroll int) string {
	clear := services.ReplaceTabs(r.l.input)
	amount := max(r.w-len(clear)-enums.InitialOffset-1, 0)
	return strings.Repeat(" ", amount+xScroll)
}

func (r *Render) RenderCodeLine(line []rune, show bool, xScroll int) (string, render.Render, bool) {
	r.l.input = line
	if string(line) == "```" {
		if !show {
			line = []rune("   ")
		}
		l := r.theme.Background + string(line) + r.fillSpace(xScroll)
		return l, render.Markdown, true
	}

	r.l.position = 0
	r.l.readPosition = 0
	r.l.readChar()

	var mode render.Render = render.Code
	var data strings.Builder

	for tok := r.l.NextToken(); tok.Type != EOL; tok = r.l.NextToken() {
		switch tok.Type {
		case keyword:
			data.WriteString(r.renderKeyWord(&tok))
		case str:
			data.WriteString(r.renderString(&tok))
		case number:
			data.WriteString(r.renderNumber(&tok))
		case bracket:
			data.WriteString(r.renderBracket(&tok))
		case comment:
			data.WriteString(r.renderComment(&tok))
		case symbol:
			data.WriteString(r.renderSymbol(&tok))
		case text:
			data.WriteString(string(tok.Literal))
		case whiteSpace:
			data.WriteString(string(tok.Literal))
		case wseol:
			data.WriteString(r.renderWSEOL(&tok))
		case tab:
			data.WriteString(r.renderTab(&tok))
		}
	}
	l := r.theme.Background + data.String() + r.fillSpace(xScroll)
	return l, mode, false
}

func (r *Render) renderTab(t *Token) string {
	cfg := config.GetConfig()
	if len(t.Literal) == cfg.TabStop {
		return r.theme.Comment + ascii.CodeTab.Str() + ascii.ResetFg.Str() + string(t.Literal[1:])
	} else {
		return r.theme.Comment + ascii.Tab.Str() + ascii.ResetFg.Str() + string(t.Literal[1:])
	}
}

func (r *Render) renderWSEOL(t *Token) string {
	return strings.Repeat(r.theme.Comment+ascii.WSEOL.Str(), len(t.Literal))
}

func (r *Render) renderBracket(t *Token) string {
	return r.theme.Bracket + string(t.Literal) + ascii.ResetFg.Str()
}

func (r *Render) renderSymbol(t *Token) string {
	return r.theme.Symbol + string(t.Literal) + ascii.ResetFg.Str()
}

func (r *Render) renderKeyWord(t *Token) string {
	return r.theme.Keyword + string(t.Literal) + ascii.ResetFg.Str()
}

func (r *Render) renderNumber(t *Token) string {
	return r.theme.Number + string(t.Literal) + ascii.ResetFg.Str()
}

func (r *Render) renderString(t *Token) string {
	return r.theme.String + string(t.Literal) + ascii.ResetFg.Str()
}

func (r *Render) renderComment(t *Token) string {
	return r.theme.Comment + string(t.Literal) + ascii.ResetFg.Str()
}
