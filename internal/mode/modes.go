package mode

// About |Mode|
// Mode is just an aliase for string
// I don't rly know why is it in separate package, but i'll keep it here for now
type Mode string

const (
	Normal      Mode = "NORMAL"
	Command     Mode = "COMMAND"
	Insert      Mode = "INSERT"
	Visual      Mode = "VISUAL"
	Visual_line Mode = "VISUAL LINE"
	Replace     Mode = "REPLACE"
)
