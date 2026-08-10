package services

func ReplaceTabs(line []rune, tabstop int) []rune {
	var new []rune

	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tabstop - (visual % tabstop)
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

func ReadTabAt(line []rune, index int, tabstop int) []rune {
	var new []rune

	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tabstop - (visual % tabstop)
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

func CursorShiftAt(line []rune, index int, tabstop int) int {

	shift := 0
	visual := 0
	for i := range index {
		if i == len(line) {
			return shift
		}
		if line[i] == '\t' {
			tab_width := tabstop - (visual % tabstop)
			visual += tab_width
			shift += tab_width - 1
		} else {
			visual += 1
		}
	}

	return shift
}

func CursorShift(line []rune, tabstop int) int {

	shift := 0
	visual := 0
	for i := range line {
		if line[i] == '\t' {
			tab_width := tabstop - (visual % tabstop)
			visual += tab_width
			shift += tab_width - 1
		} else {
			visual += 1
		}
	}

	return shift
}

func CursorShiftCalculateAt(line []rune, index int, tabstop int) int {
	shift := tabstop - (index % tabstop)
	return shift
}
