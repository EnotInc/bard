package buffers

type BufferType int

const (
	_ BufferType = iota
	Markdown
	Code
	Other
)

var CodeExt map[string]bool = map[string]bool{
	"asm":               true,
	"bash":              true,
	"c":                 true,
	"cpp":               true,
	"cs":                true,
	"css":               true,
	"dart":              true,
	"db":                true,
	"env":               true,
	"erb":               true,
	"ex":                true,
	"exs":               true,
	"elixir":            true,
	"flutter":           true,
	"go":                true,
	"hs":                true,
	"haskell":           true,
	"html":              true,
	"java":              true,
	"js":                true,
	"javascript":        true,
	"json":              true,
	"kt":                true,
	"kotlin":            true,
	"kts":               true,
	"log":               true,
	"lua":               true,
	"mk":                true,
	"Makefile":          true,
	"make":              true,
	"package.json":      true,
	"package-lock.json": true,
	"perl":              true,
	"pl":                true,
	"php":               true,
	"py":                true,
	"python":            true,
	"rs":                true,
	"rust":              true,
	"sh":                true,
	"shell":             true,
	"sql":               true,
	"swift":             true,
	"toml":              true,
	"ts":                true,
	"typescript":        true,
	"xul":               true,
	"xml":               true,
	"xhtml":             true,
	"yml":               true,
	"yaml":              true,
}
