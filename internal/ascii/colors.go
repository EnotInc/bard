package ascii

// About |Color|
// Color is just a string
// Each Color has ascii escape sequence
type Color string

func (a Color) Str() string {
	return string(a)
}

const (
	// -===[ Basic colors ]===-
	Reset    Color = "\033[0m"
	ResetFg  Color = "\033[39m"
	RedFg    Color = "\033[31m"
	GreenFg  Color = "\033[32m"
	GrayFg   Color = "\033[90m"
	YellowFg Color = "\033[33m"
	PurpleFg Color = "\033[35m"
	CyanFg   Color = "\033[36m"
	BlueFg   Color = "\033[94m"
	GrayBg   Color = "\033[100m"

	// -===[ Text style ]===-
	Bold       Color = "\033[1m"
	Italic     Color = "\033[3m"
	BoldItalic Color = "\033[1m\033[3m"
	Underline  Color = "\033[4m"
	Stricked   Color = "\033[9m"
	Hightlight Color = "\033[43m"

	// -===[ Symbols ]===-
	SymbolColor Color = GrayFg
	Quote       Color = GreenFg
	ListColor   Color = PurpleFg
	TagColor    Color = PurpleFg
	CodeLine    Color = YellowFg
	Header      Color = BlueFg
	Tab         Color = BlueFg
	WSEOLColor  Color = "\033[91m"
	CodeBg      Color = "\033[48;5;234m"
	Link        Color = "\033[4;36m"
	LowerBarBg  Color = "\033[48;5;16m"
)
