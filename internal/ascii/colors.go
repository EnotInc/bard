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
	Reset   Color = "\033[0m"
	ResetFg Color = "\033[39m"
	Error   Color = "\033[31m"

	// -===[ Text style ]===-
	Bold       Color = "\033[1m"
	Italic     Color = "\033[3m"
	BoldItalic Color = "\033[1m\033[3m"
	Stricked   Color = "\033[9m"
)
