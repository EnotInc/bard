package general

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/internal/ascii"
)

func PainAsAttr(symbol string) string {
	sym := PaintString(ascii.SymbolColor, symbol)
	return sym + ascii.Reset.Str()
}

func PaintString(c ascii.Color, str string) string {
	var s strings.Builder
	for _, x := range str {
		fmt.Fprint(&s, c, string(x))
	}
	return s.String()
}
