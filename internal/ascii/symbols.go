package ascii

// Used simply to store some specific unicode symbols
type Symbol string

func (s Symbol) Str() string {
	return string(s)
}

const (
	QuoteSymbol Symbol = "\u2503"
	Shield      Symbol = "\\"
	SplitLIne   Symbol = "\u2015"
	ListDash    Symbol = "\u2981"
	BoxEmpty    Symbol = " \u2610"
	BoxField    Symbol = " \u25a0"
	BoxComplete Symbol = " \u2612"

	WSEOL   Symbol = "\u00b7"
	NewLine Symbol = "\u21b5"
	Tab     Symbol = "\u21a6"

	TagS Symbol = "["
	TagE Symbol = "]"

	Cursor Symbol = "\u2592"
)
