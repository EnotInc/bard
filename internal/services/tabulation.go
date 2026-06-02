package services

import (
	"fmt"
	"strings"

	"github.com/EnotInc/Bard/config"
)

func ReplaceTabs(line []rune) []rune {
	var new strings.Builder
	tw := config.Get().TabStop

	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tw - (visual % tw)
			visual += tab_width

			spaces := strings.Repeat(" ", tab_width)
			fmt.Fprint(&new, spaces)
		} else {
			fmt.Fprintf(&new, "%c", line[i])
			visual += 1
		}
	}

	return []rune(new.String())
}

func ReadTabAt(line []rune, index int) []rune {
	var new strings.Builder
	tw := config.Get().TabStop

	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tw - (visual % tw)
			visual += tab_width

			if i == index {
				spaces := strings.Repeat(" ", tab_width)
				fmt.Fprint(&new, spaces)
			}
		} else {
			visual += 1
		}
	}

	return []rune(new.String())
}

func CursorShiftAt(line []rune, index int) int {
	tw := config.Get().TabStop

	shift := 0
	visual := 0
	for i := range index {
		if i == len(line) {
			return shift
		}
		if line[i] == '\t' {
			tab_width := tw - (visual % tw)
			visual += tab_width
			shift += tab_width - 1
		} else {
			visual += 1
		}
	}

	return shift
}

func CursorShift(line []rune) int {
	tw := config.Get().TabStop

	shift := 0
	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tw - (visual % tw)
			visual += tab_width
			shift += tab_width - 1
		} else {
			visual += 1
		}
	}

	return shift
}
