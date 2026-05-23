package code

type TokenType int

type Token struct {
	Type    TokenType
	Literal []rune
}

const (
	_          TokenType = iota
	keyword              // for if ...
	str                  // "string"
	number               // 1234...
	bracket              // ({[]})
	symbol               // + = - ...
	comment              // '//' or '#'
	text                 // any other text
	whiteSpace           //
	wseol                // white space at the end of line
	EOL                  // End of line
	tab
)
