package general

import (
	"fmt"

	"github.com/EnotInc/Bard/internal/ascii"
)

func PainAsAttr(symbol string) string {
	sym := PaintString(ascii.SymbolColor, symbol)
	return sym + ascii.Reset.Str()
}

func PaintString(c ascii.Color, str string) string {
	var s = ""
	for _, x := range str {
		s += fmt.Sprintf("%s%c", c, x)
	}
	return s
}
