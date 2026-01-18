package render

type TokenType string

type Token struct {
	Type    TokenType
	Literal []rune
}

const (
	TEXT = "text"

	Unknow = "unknow"
	Symbol = "symbol"

	EOL = "EOL" // End Of Line

	LineNumber = "n."
	LineDash   = "-"

	OneStar    = "*"
	TwoStars   = "**"
	ThreeStars = "***"

	OneUnderline    = "_"
	TwoUnderlines   = "__"
	ThreeUnderlines = "___"

	Stricked = "~~"

	Header_1 = "#"
	Header_2 = "##"
	Header_3 = "###"
	Header_4 = "####"
	Header_5 = "#####"
	Header_6 = "######"

	Quote = ">"
)
