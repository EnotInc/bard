package markdown

type TokenType int

type Token struct {
	Literal []rune
	Value   []rune
	Type    TokenType
}

const ( // maybe it's better to use iota for tokens?
	_    TokenType = iota
	text           // "text"
	tab

	whitespace // " "
	wseol      // "." - White Space at the End Of Line
	symbol     // "symbol"
	shield     // "\\"

	eol // "eol" - End Of Line

	listNumberDot // "n."
	listNumberB   // "n)"
	listDash      // "-"
	listBoxEmpty  // "- [ ]"
	listBoxFilled // "- [x]"

	oneStar    // "*"
	twoStars   // "**"
	threeStars // "***"

	oneUnderLine    // "_"
	twoUnderLines   // "__"
	threeUnderLines // "___"

	stricked   // "~~"
	hightlight // "=="

	header_1 // "#"
	header_2 // "##"
	header_3 // "###"
	header_4 // "####"
	header_5 // "#####"
	header_6 // "######"

	quote // ">"
	tag   // "#text"

	codeLine  // "`"
	codeBlock // "```"

	link  // "[text](link)"
	image // "![text](link)"

	html // any html block in <>
)
