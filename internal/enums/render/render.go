package enums

// helps to decide which type of render use
type Render int

const (
	_ Render = iota
	Markdown
	Code
)
