package ascii

const (
	ClearView    = "\033[2J"
	ClearHistory = "\033[3J"
	MoveToStart  = "\033[0H"
	CursorReset  = "\033]112\a"

	SaveTerminal  = "\033[?1049h"
	ResetTerminal = "\033[?1049l"

	CursorBloc      = "\x1b[2 q"
	CursorLine      = "\x1b[6 q"
	CursorUnderline = "\x1b[4 q"
)
