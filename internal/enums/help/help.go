package help

// Used to work with `:h <topic>` command
type Topic string

const (
	About      Topic = "about"
	Modes      Topic = "modes"
	Normal     Topic = "normal"
	Insert     Topic = "insert"
	Command    Topic = "command"
	Visual     Topic = "visual"
	VisualLine Topic = "visual-line"
	Config     Topic = "config"
	Explorer   Topic = "explorer"
	Space      Topic = "spase"
	Theme      Topic = "theme"
)
