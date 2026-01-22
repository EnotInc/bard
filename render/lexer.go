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
	case '\\':
		if isLetterOrNumber(l.peekChar()) || l.peekChar() == 0 || l.peekChar() == ' ' {
			t = Token{Type: Symbol, Literal: []rune{l.ch}}
		} else {
			t = Token{Type: Shield, Literal: []rune{l.ch}}
		}
		l.readChar()
	case '>':
		t = Token{Type: Quote, Literal: []rune{'>'}}
		l.readChar()
	case ' ':
		l.readChar()
		t = Token{Type: WhiteSpace, Literal: []rune{' '}}
	case '*':
		t = l.getAttrToken('*', []TokenType{OneStar, TwoStars, ThreeStars})
	case '_':
		t = l.getAttrToken('_', []TokenType{OneUnderline, TwoUnderlines, ThreeUnderlines})
	case '~':
		if l.peekChar() == '~' {
			l.readChar()
			t = Token{Type: Stricked, Literal: []rune("~~")}
		} else {
			t = Token{Type: Symbol, Literal: []rune{l.ch}}
		}
		l.readChar()
	case '#':
		pos := l.position
		count := 1

		for l.peekChar() == '#' {
			count += 1
			l.readChar()
		}

		end := l.position + 1
		lit := []rune(l.input[pos:end])

		if count == 1 && isLetterOrNumber(l.peekChar()) /* && l.peekChar() != ' '*/ {
			l.readChar()
			text := l.readText()
			t = Token{Type: Tag, Literal: append(lit, text...)}
		} else if count > 6 || l.peekChar() != ' ' {
			t = Token{Type: Symbol, Literal: lit}
			l.readChar()
		} else {
			switch count {
			case 1:
				t = Token{Type: Header_1, Literal: lit}
			case 2:
				t = Token{Type: Header_2, Literal: lit}
			case 3:
				t = Token{Type: Header_3, Literal: lit}
			case 4:
				t = Token{Type: Header_4, Literal: lit}
			case 5:
				t = Token{Type: Header_5, Literal: lit}
			case 6:
				t = Token{Type: Header_6, Literal: lit}
			}
			l.readChar()
		}
	case '`':
		pos := l.position
		count := 1

		for l.peekChar() == '`' {
			count += 1
			l.readChar()
		}

		end := l.position + 1
		lit := []rune(l.input[pos:end])

		switch count {
		case 1:
			t = Token{Type: CodeLine, Literal: lit}
		case 3:
			t = Token{Type: CodeBlock, Literal: lit}
		default:
			t = Token{Type: Symbol, Literal: lit}
		}
		l.readChar()
	case 0:
		t = Token{Type: EOL, Literal: []rune("")}
	default:
		if isLetterOrNumber(l.ch) {
			s := l.readText()
			t = Token{Type: TEXT, Literal: s}
		} else {
			t = Token{Type: Symbol, Literal: []rune{l.ch}}
			l.readChar()
		}
	}

	return t
}

func (l *lexer) readText() []rune {
	pos := l.position
	for isLetterOrNumber(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isLetterOrNumber(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9'
}

func (l *lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *lexer) getAttrToken(ch rune, types []TokenType) Token {
	var t Token
	pos := l.position
	count := 1

	for l.peekChar() == ch {
		count += 1
		l.readChar()
	}

	end := l.position + 1
	if count > 3 || l.peekChar() == ' ' {
		t = Token{Type: Symbol, Literal: []rune(l.input[pos:end])}
	} else {
		switch count {
		case 1:
			t = Token{Type: types[0], Literal: []rune(l.input[pos:end])}
		case 2:
			t = Token{Type: types[1], Literal: []rune(l.input[pos:end])}
		case 3:
			t = Token{Type: types[2], Literal: []rune(l.input[pos:end])}
		}
	}

	l.readChar()
	return t
}
