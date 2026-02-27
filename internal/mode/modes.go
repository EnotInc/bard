package mode

type Mode string

const (
	Normal      Mode = "NORMAL"
	Command     Mode = "COMMAND"
	Insert      Mode = "INSERT"
	Visual      Mode = "VISUAL"
	Visual_line Mode = "VISUAL LINE"
	Replace     Mode = "REPLACE"
)
