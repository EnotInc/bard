package code

import (
	"strings"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
	"github.com/EnotInc/Bard/internal/render/general"
)

type Render struct {
	l *Lexer
	w int
}

func NewRender(w int) *Render {
	r := &Render{w: w}
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

func (r *Render) fillSpace() string {
	amount := max(r.w-len(r.l.input)-enums.InitialOffset-1, 0)
	return strings.Repeat(" ", amount)
}

func (r *Render) RenderCodeLine(line []rune, show bool) (string, int, enums.Render) {
	r.l.input = line
	if string(line) == "```" {
		if !show {
			line = []rune("   ")
		}
		l := ascii.CodeBg.Str() + string(line) + r.fillSpace()
		diff := -r.w
		return l, diff, enums.Markdown
	}

	r.l.position = 0
	r.l.readPosition = 0
	r.l.readChar()

	var mode enums.Render = enums.Code
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
		}
	}
	l := ascii.CodeBg.Str() + data.String() + r.fillSpace()
	diff := -r.w - len(r.l.input) - enums.InitialOffset - 1
	return l, diff, mode
}

func (r *Render) renderWSEOL(t *Token) string {
	return strings.Repeat(ascii.WSEOLColor.Str()+ascii.WSEOL.Str(), len(t.Literal))
}

func (r *Render) renderBracket(t *Token) string {
	return general.PaintString(ascii.PurpleFg, string(t.Literal)) + ascii.ResetFg.Str()
}

func (r *Render) renderSymbol(t *Token) string {
	return general.PaintString(ascii.YellowFg, string(t.Literal)) + ascii.ResetFg.Str()
}

func (r *Render) renderKeyWord(t *Token) string {
	return general.PaintString(ascii.YellowFg, string(t.Literal)) + ascii.ResetFg.Str()
}

func (r *Render) renderNumber(t *Token) string {
	return general.PaintString(ascii.PurpleFg, string(t.Literal)) + ascii.ResetFg.Str()
}

func (r *Render) renderString(t *Token) string {
	return general.PaintString(ascii.GreenFg, string(t.Literal)) + ascii.ResetFg.Str()
}

func (r *Render) renderComment(t *Token) string {
	return general.PaintString(ascii.GrayFg, string(t.Literal)) + ascii.ResetFg.Str()
}
