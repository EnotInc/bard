package enums

// helps to decide wich type of render use
type Render int

const (
	Raw Render = iota
	Markdown
	Code
)
