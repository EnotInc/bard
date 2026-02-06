package render

type lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
}

func NewLexer() *lexer {
	//l := &lexer{input: input}
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
	case '-':
		//NOTE: I know that this is not looking good. I'll figure this out, later, maybe...
		str := "-"
		if l.peekChar() == ' ' {
			l.readChar()
			str += " "
			isField := false
			t = Token{Type: ListDash, Literal: []rune{l.ch}}
			if l.peekChar() == '[' {
				l.readChar()
				str += "["
				if l.peekChar() == ' ' || isLetter(l.peekChar()) || l.peekChar() == '?' {
					isField = ' ' != l.peekChar()
					l.readChar()
					str += string(l.ch)
					if l.peekChar() == ']' {
						str += "]"
						l.readChar()
						if isField {
							t = Token{Type: ListBoxField, Literal: []rune(str)}
						} else {
							t = Token{Type: ListBoxEmpty, Literal: []rune(str)}
						}
					} else {
						t = Token{Type: TEXT, Value: []rune(str)}
					}
				} else {
					t = Token{Type: TEXT, Value: []rune(str)}
				}
			} else {
				t = Token{Type: ListDash, Literal: []rune(str)}
			}
		} else {
			t = Token{Type: Symbol, Value: []rune(str)}
		}
		l.readChar()
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
	if count > 3 /*|| l.peekChar() == ' '*/ {
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
