package code

import (
	"github.com/EnotInc/Bard/config"
	txt "github.com/EnotInc/Bard/internal/services/text"
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

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readNumber() []rune {
	pos := l.position
	for txt.IsNumber(l.ch) || (l.ch == '.' && txt.IsNumber(l.peekChar())) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readText() []rune {
	pos := l.position
	for txt.IsLetterOrNumber(l.ch) || islinkSymbol(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readString(q rune) []rune {
	pos := l.position
	l.readChar()
	for l.ch != q && l.ch != 0 {
		l.readChar()
	}
	if l.ch != 0 {
		l.readChar()
	}
	str := l.input[pos:l.position]
	return str
}

func (l *Lexer) readWhiteSpace() ([]rune, bool) {
	pos := l.position
	for l.ch == ' ' && l.ch != 0 {
		l.readChar()
	}
	return l.input[pos:l.position], l.ch != 0
}

func islinkSymbol(ch rune) bool {
	return ch == '/' || ch == '.' || ch == ':' || ch == '?' || ch == '=' || ch == '%'
}

func (l *Lexer) readTab() Token {
	ts := config.GetConfig().TabStop
	new := txt.ReadTabAt(l.input, l.position, ts)
	return Token{Type: tab, Literal: new}
}

func (l *Lexer) NextToken() Token {
	var t Token
	switch l.ch {
	case '/', '-':
		pos := l.position
		c := l.ch
		l.readChar()
		if l.ch == c {
			s := l.input[pos:]
			t = Token{Type: comment, Literal: s}
			l.readPosition = len(l.input)
			l.readChar()
		} else {
			t = Token{Type: symbol, Literal: []rune{c}}
		}
	case '\t':
		t = l.readTab()
		l.readChar()
	case '#':
		s := l.input[l.position:]
		t = Token{Type: comment, Literal: s}
		l.readPosition = len(l.input)
		l.readChar()
	case '(', '[', '{', '}', ']', ')':
		t = Token{Type: bracket, Literal: []rune{l.ch}}
		l.readChar()
	case '"', '\'':
		s := l.readString(l.ch)
		t = Token{Type: str, Literal: s}
	case ' ':
		spaces, isEnd := l.readWhiteSpace()
		if isEnd {
			t = Token{Type: whiteSpace, Literal: spaces}
		} else {
			t = Token{Type: wseol, Literal: spaces}
		}
	case 0:
		t = Token{Type: EOL, Literal: []rune("")}
		l.readChar()
	default:
		if txt.IsNumber(l.ch) {
			s := l.readNumber()
			t = Token{Type: number, Literal: s}
		} else if txt.IsLetter(l.ch) {
			s := l.readText()
			if _, ok := keywords[string(s)]; ok {
				t = Token{Type: keyword, Literal: s}
			} else {
				t = Token{Type: text, Literal: s}
			}
			//l.readChar()
		} else {
			t = Token{Type: symbol, Literal: []rune{l.ch}}
			l.readChar()
		}
	}
	return t
}
