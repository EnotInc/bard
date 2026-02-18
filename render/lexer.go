package render

type lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
}

func NewLexer() *lexer {
	l := &lexer{}
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
	case '[':
		t = l.readLink()
	case '!':
		if l.peekChar() == '[' {
			l.readChar()
			t = l.readLink()
			if t.Type == Symbol {
				t.Value = append([]rune{'!'}, t.Value...)
			} else {
				t.Type = Image
				t.Literal = append([]rune{'!'}, t.Literal...)
				t.Value = append([]rune{'!'}, t.Value...)
			}
		} else {
			t = Token{Type: Symbol, Value: []rune("!")}
			l.readChar()
		}
	case '-':
		t = l.readListOrCheckBox()
	case '\\':
		if isNumber(l.peekChar()) || isLetter(l.peekChar()) || l.peekChar() == 0 || l.peekChar() == ' ' {
			t = Token{Type: Symbol, Value: []rune{l.ch}}
		} else {
			sh := l.ch
			l.readChar()
			t = Token{Type: Shield, Literal: []rune{sh}, Value: []rune{l.ch}}
		}
		l.readChar()
	case '>':
		t = Token{Type: Quote, Literal: []rune{l.ch}}
		l.readChar()
	case ' ':
		l.readChar()
		t = Token{Type: WhiteSpace, Value: []rune{' '}}
	case '*':
		t = l.getAttrToken('*', []TokenType{OneStar, TwoStars, ThreeStars})
	case '_':
		t = l.getAttrToken('_', []TokenType{OneUnderline, TwoUnderlines, ThreeUnderlines})
	case '~':
		if l.peekChar() == '~' {
			l.readChar()
			t = Token{Type: Stricked, Literal: []rune("~~")}
		} else {
			t = Token{Type: Symbol, Value: []rune{l.ch}}
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

		if count == 1 && (isLetter(l.peekChar()) || isNumber(l.peekChar())) {
			l.readChar()
			text := l.readText()
			t = Token{Type: Tag, Literal: lit, Value: text}
		} else if count > 6 || l.peekChar() != ' ' {
			t = Token{Type: Symbol, Value: lit}
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
		// case 3:
		// 	t = Token{Type: CodeBlock, Literal: lit}
		default:
			t = Token{Type: Symbol, Value: lit}
		}
		l.readChar()
	case 0:
		t = Token{Type: EOL, Literal: []rune("")}
	default:
		if isNumber(l.ch) {
			s := l.readNumber()
			switch l.ch {
			case ')':
				t = Token{Type: ListNumberB, Value: s, Literal: []rune{')'}}
				l.readChar()
			case '.':
				t = Token{Type: ListNumberDot, Value: s, Literal: []rune{'.'}}
				l.readChar()
			default:
				t = Token{Type: TEXT, Value: s}
			}
		} else if isLetter(l.ch) || isNumber(l.ch) {
			s := l.readText()
			t = Token{Type: TEXT, Value: s}
		} else {
			t = Token{Type: Symbol, Value: []rune{l.ch}}
			l.readChar()
		}
	}

	return t
}

func (l *lexer) readNumber() []rune {
	pos := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *lexer) readText() []rune {
	pos := l.position
	for isLetter(l.ch) || isNumber(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumber(ch rune) bool {
	return '0' <= ch && ch <= '9'
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
	if count > 3 {
		t = Token{Type: Symbol, Value: []rune(l.input[pos:end])}
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

func (l *lexer) readLink() Token {
	text := l.position
	l.readChar()
	start := l.position
	for l.ch != ']' && l.ch != 0 {
		if l.ch == '\\' && l.peekChar() == ']' {
			l.readChar()
			l.readChar()
			continue
		}
		l.readChar()
	}
	if l.ch == ']' {
		l.readChar()

		if l.ch == '(' {
			txt := l.input[text+1 : l.position-1]
			l.readChar()

			for l.ch != ')' && l.ch != 0 {
				if l.ch == '\\' && l.peekChar() == ')' {
					l.readChar()
					l.readChar()
					continue
				}
				l.readChar()
			}
			if l.ch == ')' {
				l.readChar()
				lnk := l.input[text:l.position]
				return Token{Type: Link, Value: txt, Literal: lnk}
			}
		}
	}

	l.position = start
	l.readPosition = start
	l.readChar()
	return Token{Type: Symbol, Value: []rune("[")}
}

func (l *lexer) readListOrCheckBox() Token {
	pos := l.position
	if l.peekChar() != ' ' {
		l.readChar()
		return Token{Type: Symbol, Value: []rune{'-'}}
	}

	l.readChar()
	start := l.position

	if l.peekChar() == '[' {
		l.readChar()

		ch := l.peekChar()
		if ch == ' ' || isLetter(ch) || ch == '?' {
			l.readChar()
			isField := ch != ' '

			if l.peekChar() == ']' {
				l.readChar() // reading the ']' symbol
				l.readChar() // reading next symbol for lexer
				literal := l.input[pos:l.position]

				if isField {
					return Token{Type: ListBoxField, Literal: literal}
				} else {
					return Token{Type: ListBoxEmpty, Literal: literal}
				}
			}
		}

		l.position = start
		l.readPosition = start
		l.readChar()
		return Token{Type: ListDash, Literal: []rune("-")}
	}

	return Token{Type: ListDash, Literal: []rune("-")}
}
