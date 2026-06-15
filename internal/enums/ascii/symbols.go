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
)
