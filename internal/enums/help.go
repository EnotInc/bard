package enums

// About |Help|
// Used to work with `:h <topic>` command
type Help string

const (
	About   Help = "about"
	Modes   Help = "modes"
	Normal  Help = "normal"
	Insert  Help = "insert"
	Command Help = "command"
	Visual  Help = "visual"
)
