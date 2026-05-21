package general

import (
	"fmt"
	"strings"
)

func PaintString(color string, str string) string {
	var s strings.Builder
	for _, x := range str {
		fmt.Fprint(&s, color, string(x))
	}
	return s.String()
}
