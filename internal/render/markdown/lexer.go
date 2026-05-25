package markdown

import (
	"github.com/EnotInc/Bard/internal/services"
)

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
}

func newLexer() *Lexer {
	l := &Lexer{}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var t Token

	switch l.ch {
	case '\t':
		t = l.readTab()
		l.readChar()
	case '[':
		t = l.readLink()
	case '!':
		if l.peekChar() == '[' {
			l.readChar()
			t = l.readLink()
			if t.Type == symbol {
				t.Value = append([]rune{'!'}, t.Value...)
			} else {
				t.Type = image
				t.Literal = append([]rune{'!'}, t.Literal...)
				t.Value = append([]rune{'!'}, t.Value...)
			}
		} else {
			t = Token{Type: symbol, Value: []rune("!")}
			l.readChar()
		}
	case '-':
		t = l.readListOrCheckBox()
	case '\\':
		if (isNumber(l.peekChar()) || isLetter(l.peekChar()) || l.peekChar() == 0 || l.peekChar() == ' ') && l.peekChar() != '_' {
			t = Token{Type: symbol, Value: []rune{l.ch}}
		} else {
			sh := l.ch
			l.readChar()
			t = Token{Type: shield, Literal: []rune{sh}, Value: []rune{l.ch}}
		}
		l.readChar()
	case '>':
		t = Token{Type: quote, Literal: []rune{l.ch}}
		l.readChar()
	case ' ':
		spaces, isEnd := l.readWhiteSpace()
		if isEnd {
			t = Token{Type: whitespace, Value: spaces}
		} else {
			t = Token{Type: wseol, Value: spaces}
		}
	case '*':
		t = l.getAttrToken('*', []TokenType{oneStar, twoStars, threeStars})
	case '_':
		t = l.getAttrToken('_', []TokenType{oneUnderLine, twoUnderLines, threeUnderLines})
	case '~':
		if l.peekChar() == '~' {
			l.readChar()
			t = Token{Type: stricked, Literal: []rune("~~")}
		} else {
			t = Token{Type: symbol, Value: []rune{l.ch}}
		}
		l.readChar()
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			t = Token{Type: hightlight, Value: []rune("==")}
		} else {
			t = Token{Type: symbol, Value: []rune("=")}
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
			t = Token{Type: tag, Literal: lit, Value: text}
		} else if count > 6 || l.peekChar() != ' ' {
			t = Token{Type: symbol, Value: lit}
			l.readChar()
		} else {
			switch count {
			case 1:
				t = Token{Type: header_1, Literal: lit}
			case 2:
				t = Token{Type: header_2, Literal: lit}
			case 3:
				t = Token{Type: header_3, Literal: lit}
			case 4:
				t = Token{Type: header_4, Literal: lit}
			case 5:
				t = Token{Type: header_5, Literal: lit}
			case 6:
				t = Token{Type: header_6, Literal: lit}
			}
			l.readChar()
		}
	case '`':
		pos := l.position
		l.readChar()
		count := 1
		for l.ch == '`' {
			l.readChar()
			count += 1
		}
		lit := l.input[pos:l.position]

		switch count {
		case 1:
			s := l.readCodeLine()
			t = Token{Type: codeLine, Literal: []rune{'`'}, Value: s}
		case 2:
			t = Token{Type: symbol, Value: lit}
		case 3:
			t = Token{Type: codeBlock, Literal: lit, Value: l.input[l.position:]}
			l.position = len(l.input)
			l.readPosition = len(l.input)
		default:
			t = Token{Type: symbol, Value: lit}
		}
	case '<':
		t = l.readHTMLBlock()
	case 0:
		t = Token{Type: eol, Literal: []rune("")}
	default:
		if isNumber(l.ch) {
			s := l.readNumber()
			switch l.ch {
			case ')':
				t = Token{Type: listNumberB, Value: s, Literal: []rune{')'}}
				l.readChar()
			case '.':
				t = Token{Type: listNumberDot, Value: s, Literal: []rune{'.'}}
				l.readChar()
			default:
				t = Token{Type: text, Value: s}
			}
		} else if isLetter(l.ch) || isNumber(l.ch) {
			s := l.readText()
			t = Token{Type: text, Value: s}
		} else {
			t = Token{Type: symbol, Value: []rune{l.ch}}
			l.readChar()
		}
	}

	return t
}

func (l *Lexer) readTab() Token {
	new := services.ReadTabAt(l.input, l.position)
	return Token{Type: tab, Literal: []rune(new)}
}

func (l *Lexer) readNumber() []rune {
	pos := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readText() []rune {
	pos := l.position
	for isLetter(l.ch) || isNumber(l.ch) {
		l.readChar()
		if l.ch == '_' && (l.peekChar() == '_' || l.peekChar() == ' ' || l.peekChar() == 0) {
			break
		}
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readHTMLBlock() Token {
	start := l.position
	l.readChar()
	for l.ch != '>' && l.peekChar() != 0 && (isLetter(l.ch) || isNumber(l.ch) || isSymbol(l.ch) || l.ch == ' ') {
		l.readChar()
	}

	value := l.input[start:l.position]

	if len(value) == 1 {
		return Token{
			Type:  symbol,
			Value: value,
		}
	}

	var literal []rune

	if len(value) > 2 && value[0] == '<' && value[1] == '/' && value[len(value)-1] == '>' { // </>
		literal = []rune{value[0], value[1], value[len(value)-1]}
		value = value[2 : len(value)-1]

	} else if len(value) >= 2 && value[0] == '<' && value[len(value)-1] == '>' { // <>
		literal = []rune{value[0], value[len(value)-1]}
		value = value[1 : len(value)-1]

	} else if len(value) >= 2 && value[0] == '<' && value[1] == '/' { // </
		literal = []rune{value[0], value[1]}
		value = value[2:]

	} else { // <
		literal = []rune{'<'}
		if len(value) > 1 {
			value = value[1:]
		} else {
			value = []rune("")
		}
	}

	return Token{
		Type:    html,
		Value:   value,
		Literal: literal,
	}
}

func (l *Lexer) readCodeLine() []rune {
	pos := l.position
	l.readChar()
	for l.ch != '`' && l.ch != 0 {
		l.readChar()
	}
	l.readChar()
	return l.input[pos:l.position]
}

// TODO: move to services 'is letter or number' functions
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumber(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isSymbol(ch rune) bool {
	return ch == '+' || ch == '/' || ch == '\\'
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) getAttrToken(ch rune, types []TokenType) Token {
	var t Token
	pos := l.position
	count := 1

	for l.peekChar() == ch {
		count += 1
		l.readChar()
	}

	end := l.position + 1
	if count > 3 {
		t = Token{Type: symbol, Value: []rune(l.input[pos:end])}
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

func (l *Lexer) readWhiteSpace() ([]rune, bool) {
	pos := l.position
	for l.ch == ' ' && l.ch != 0 {
		l.readChar()
	}
	return l.input[pos:l.position], l.ch != 0
}

func (l *Lexer) readLink() Token {
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
				return Token{Type: link, Value: txt, Literal: lnk}
			}
		}
	}

	l.position = start
	l.readPosition = start
	l.readChar()
	return Token{Type: symbol, Value: []rune("[")}
}

func (l *Lexer) readListOrCheckBox() Token {
	pos := l.position
	if l.peekChar() != ' ' {
		l.readChar()
		return Token{Type: symbol, Value: []rune{'-'}}
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
					return Token{Type: listBoxFilled, Literal: literal}
				} else {
					return Token{Type: listBoxEmpty, Literal: literal}
				}
			}
		}

		l.position = start
		l.readPosition = start
		l.readChar()
		return Token{Type: listDash, Literal: []rune("-")}
	}

	return Token{Type: listDash, Literal: []rune("-")}
}
