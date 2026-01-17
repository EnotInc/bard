package render

type lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
}

func NewLexer(input []rune) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *lexer) NextToken() Token {
	var t Token
	switch l.ch {
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			if l.peekChar() == '*' {
				l.readChar()
				t = Token{Type: ThreeStars, Literal: []rune("***")}
			} else {
				t = Token{Type: TwoStars, Literal: []rune("**")}
			}
		} else {
			t = Token{Type: OneStar, Literal: []rune("*")}
		}
		l.readChar()
	case '~':
		if l.peekChar() == '~' {
			l.readChar()
			t = Token{Type: Stricked, Literal: []rune("~~")}
		} else {
			t = l.readText()
		}
	case 0:
		t = Token{Type: EOL, Literal: []rune("")}
	default:
		t = l.readText()
	}

	return t
}

func (l *lexer) readText() Token {
	pos := l.position
	for isLetterOrNumber(l.ch) {
		l.readChar()
	}

	return Token{Type: TEXT, Literal: l.input[pos:l.position]}
}

func isLetterOrNumber(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' || ch == ' '
}

func (l *lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
