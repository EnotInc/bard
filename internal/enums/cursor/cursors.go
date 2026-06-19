package cursor

type CursorType string

const (
	CursorBloc      CursorType = "\x1b[2 q"
	CursorLine      CursorType = "\x1b[6 q"
	CursorUnderline CursorType = "\x1b[4 q"
)
