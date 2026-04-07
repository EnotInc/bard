package code

type TokenType int

type Token struct {
	Type    TokenType
	Literal []rune
}

const (
	_          = iota
	keyword    // for if ...
	str        // "string"
	number     // 1234...
	bracket    // ({[]})
	symbol     // + = - ...
	comment    // '//' or '#'
	text       // any other text
	whiteSpace //
	EOL        // End of line
)
