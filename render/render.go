package render

const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	italic    = "\033[3m"
	underline = "\033[7m"
	stricked  = "\033[9m"
)

type Renderer struct {
	curMode string
	l       *lexer
}

func InitReder() *Renderer {
	r := &Renderer{}
	return r
}

func (r *Renderer) RednerMarkdownLine(line []rune, isCur bool) string {
	r.l = NewLexer(line)
	var data = ""

	for tok := r.l.NextToken(); tok.Type != EOL; tok = r.l.NextToken() {
		switch tok.Type {
		case TEXT:
			data += string(tok.Literal)
		case OneStar, OneUnderline:
			data += r.renderItalc(&tok, isCur)
		case TwoStars, TwoUnderlines:
			data += r.renderBold(&tok, isCur)
		}
	}

	data += reset
	r.curMode = reset
	return data
}

func (r *Renderer) renderHeader(t *Token) string {
	return ""
}

func (r *Renderer) renderItalc(t *Token, show bool) string {
	var s = ""
	r.changeMode(italic)
	s += r.curMode
	if show {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) renderBold(t *Token, show bool) string {
	var s = ""
	//TODO: Figure out why it's now working in windws powershell
	r.changeMode(bold)
	s += r.curMode
	if show {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) renderStreaced(t *Token, show bool) string {
	var s = ""
	r.changeMode(stricked)
	s += r.curMode
	if show {
		s += string(t.Literal)
	}

	return s
}

func (r *Renderer) changeMode(mode string) {
	if r.curMode == mode {
		r.curMode = reset
	} else {
		r.curMode = mode
	}
}
