package editor

func (e *Editor) caseVisual(key rune) {
	switch key {
	case '\033':
		e.curMode = normal
	case 'h':
		e.b.H()
	case 'j':
		e.b.J()
		e.ScrollDown()
	case 'k':
		e.b.K()
		e.ScrollUp()
	case 'l':
		e.b.L()
	}
}
