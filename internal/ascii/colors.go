package ascii

type Color string

func (a Color) Str() string {
	return string(a)
}

const (
	Reset    Color = "\033[0m"
	ResetFg  Color = "\033[39"
	RedFg    Color = "\033[31m"
	GrayFg   Color = "\033[90m"
	YellowFg Color = "\033[33m"
	CyanFg   Color = "\033[36m"
	StartSel Color = "\033[100m"

	Bold       Color = "\033[1m"
	Italic     Color = "\033[3m"
	BoldItalic Color = "\033[1m\033[3m"
	Underline  Color = "\033[4m"
	Stricked   Color = "\033[9m"

	SymbolColor Color = "\033[90m"
	Quote       Color = "\033[32m"
	CodeLine    Color = "\033[33m"
	Header      Color = "\033[94m"
	Link        Color = "\033[4;36m"
	ListColor   Color = "\033[35m"
	TagColor    Color = "\033[35m"

	StatusBar Color = "\033[100m"
	Tab       Color = "\033[94m"
)
