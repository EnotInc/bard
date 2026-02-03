package render

type TokenType string

type Token struct {
	Type    TokenType
	Literal []rune
	Value   []rune
}

const (
	TEXT = "text"

	WhiteSpace = " "
	Symbol     = "symbol"
	Shield     = "\\"

	EOL = "EOL" // End Of Line

	ListNumberDot = "n."
	ListNumberB   = "n)"
	ListDash      = "-"
	ListBoxEmpty  = "- [ ]"
	ListBoxField  = "- [x]"

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
	Tag   = "#text"

	CodeLine  = "`"
	CodeBlock = "```"
)
