package ascii

// Used simply to store some specific unicode symbols
type Symbol string

func (s Symbol) Str() string {
	return string(s)
}

const (
	QuoteSymbol Symbol = "\u2503"
	Shield      Symbol = "\\"
	SplitLine   Symbol = "\u2015"
	ListDash    Symbol = "\u2981"
	BoxEmpty    Symbol = " \u2610"
	BoxField    Symbol = " \u25a0"
	BoxComplete Symbol = " \u2612"
	LinkSymbol  Symbol = "\U0001f517"
	ImageSymbol Symbol = "\U0001f3a8"
	ColorBox    Symbol = "\u2b1b"

	QuoteNote      Symbol = "\033[34m\U0001f6c8"
	QuoteTip       Symbol = "\033[32m\U0001f4a1"
	QuoteImportant Symbol = "\033[35m\u203c"
	QuoteWarning   Symbol = "\033[33m\u26a0"
	QuoteCaution   Symbol = "\033[31m\U0001f6d1"

	WSEOL   Symbol = "\u00b7"
	NewLine Symbol = "\u21b5"
	Tab     Symbol = "\u21a6"
	CodeTab Symbol = "\u2502"

	TagS Symbol = "["
	TagE Symbol = "]"

	Cursor Symbol = "\u2592"

	BorderCUL string = "\u256d"
	BorderCUR string = "\u256e"
	BorderCDR string = "\u256f"
	BorderCDL string = "\u2570"
	BorderV   string = "\u2502"
	BorderH   string = "\u2500"

	ArrowUp   Symbol = "\u2303"
	ArrowDown Symbol = "\u2304"

	File   Symbol = "\U0001f5ce "
	Folder Symbol = "\U0001f5c0 "
	Search Symbol = "\U0001f50d"
)
