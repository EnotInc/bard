package services

import (
	"github.com/EnotInc/Bard/config"
)

func ReplaceTabs(line []rune) []rune {
	var new []rune
	tw := config.GetConfig().TabStop

	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tw - (visual % tw)
			visual += tab_width

			for range tab_width {
				new = append(new, ' ')
			}
		} else {
			new = append(new, line[i])
			visual += 1
		}
	}

	return []rune(new)
}

func ReadTabAt(line []rune, index int) []rune {
	var new []rune
	tw := config.GetConfig().TabStop

	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tw - (visual % tw)
			visual += tab_width

			if i == index {
				for range tab_width {
					new = append(new, ' ')
				}
			}
		} else {
			visual += 1
		}
	}

	return []rune(new)
}

func CursorShiftAt(line []rune, index int) int {
	tw := config.GetConfig().TabStop

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
	tw := config.GetConfig().TabStop

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

func CursorShiftCalculateAt(line []rune, index int) int {
	tw := config.GetConfig().TabStop
	shift := tw - (index % tw)
	return shift
}
