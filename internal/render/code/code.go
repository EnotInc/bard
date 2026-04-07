package code

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/ascii"
	"github.com/EnotInc/Bard/internal/enums"
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

func PaintString(c ascii.Color, str []rune) string {
	var s strings.Builder
	for _, x := range str {
		fmt.Fprint(&s, c, string(x))
	}
	return s.String() + ascii.ResetFg.Str()
}

func (r *Render) fillSpace() string {
	return strings.Repeat(" ", r.w-len(r.l.input)-enums.InitialOffset-1)
}

func (r *Render) RenderCodeLine(line []rune) (string, int, enums.Render) {
	r.l.input = line
	if string(line) == "```" {
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
			data.WriteString(" ")
		}
	}
	l := ascii.CodeBg.Str() + data.String() + r.fillSpace()
	diff := -r.w - len(r.l.input)
	return l, diff, mode
}

func (r *Render) renderBracket(t *Token) string {
	return PaintString(ascii.PurpleFg, t.Literal)
}

func (r *Render) renderSymbol(t *Token) string {
	return PaintString(ascii.YellowFg, t.Literal)
}

func (r *Render) renderKeyWord(t *Token) string {
	return PaintString(ascii.YellowFg, t.Literal)
}

func (r *Render) renderNumber(t *Token) string {
	return PaintString(ascii.PurpleFg, t.Literal)
}

func (r *Render) renderString(t *Token) string {
	return PaintString(ascii.GreenFg, t.Literal)
}

func (r *Render) renderComment(t *Token) string {
	return PaintString(ascii.GrayFg, t.Literal)
}
