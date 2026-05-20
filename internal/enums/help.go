package enums

// About |Help|
// Used to work with `:h <topic>` command
type Help string

const (
	HelpAbout      Help = "about"
	HelpModes      Help = "modes"
	HelpNormal     Help = "normal"
	HelpInsert     Help = "insert"
	HelpCommand    Help = "command"
	HelpVisual     Help = "visual"
	HelpVisualLine Help = "visual-line"
)
