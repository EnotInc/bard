package ascii

// List of ascii escape sequences, that used across whole project to work wiht termilan
const (
	ClearView    = "\033[2J"
	ClearHistory = "\033[3J"
	MoveToStart  = "\033[0H"
	CursorReset  = "\033]112\a"

	SaveTerminal  = "\033[?1049h"
	ResetTerminal = "\033[?1049l"

	ResetCursor     = "\x1b[0 q"
	CursorBloc      = "\x1b[2 q"
	CursorLine      = "\x1b[6 q"
	CursorUnderline = "\x1b[4 q"

	HideCursor = "\x1b[?25l"
	ShowCursor = "\x1b[?25h"
)
