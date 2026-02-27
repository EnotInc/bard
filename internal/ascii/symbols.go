package ascii

type Symbol string

func (s Symbol) Str() string {
	return string(s)
}

const (
	QuoteSymbol Symbol = "\u2503"
	Shield      Symbol = "\\"

	ListDash Symbol = "\u2981"
	BoxEmpty Symbol = " \u25a1"
	BoxField Symbol = " \u25a0"

	TagS Symbol = "["
	TagE Symbol = "]"
)
