package services

import (
	"strings"
)

func GetFileIcon(s string) string {
	parts := strings.Split(s, ".")
	var ext string = ""
	if len(parts) > 1 {
		ext = parts[len(parts)-1]
	} else {
		ext = strings.TrimPrefix(s, ".")
	}

	ext = strings.ToLower(ext)

	if i, ok := langIcon[ext]; ok {
		return i
	} else {
		return defaultFile
	}
}

const defaultFile = " "

var langIcon map[string]string = map[string]string{
	"asm":               "\033[34m ",
	"bash":              " ",
	"c":                 "\033[34m ",
	"config":            " ",
	"conf":              " ",
	"cfg":               " ",
	"cpp":               "\033[34m ",
	"c++":               "\033[34m ",
	"cs":                "\033[35m ",
	"c#":                "\033[35m ",
	"css":               "\033[36m ",
	"dart":              "\033[36m ",
	"db":                "\033[36m ",
	"env":               "󰳍 ",
	"erb":               "\033[31m ",
	"ex":                "\033[35m ",
	"exs":               "\033[35m ",
	"elixir":            "\033[35m ",
	"flutter":           "\033[34m ",
	"go":                "\033[36m ",
	"golang":            "\033[36m ",
	"hs":                "\033[35m ",
	"haskell":           "\033[35m ",
	"html":              "\033[33m ",
	"java":              "\033[36m ",
	"js":                "\033[33m ",
	"javascript":        "\033[33m ",
	"json":              " ",
	"kt":                "\033[34m ",
	"kotlin":            "\033[34m ",
	"kts":               "\033[34m ",
	"log":               " ",
	"lua":               "\033[34m ",
	"license":           "\033[33m ",
	"md":                "\033[33m ",
	"markdown":          "\033[33m ",
	"mk":                " ",
	"Makefile":          " ",
	"make":              " ",
	"package.json":      "\033[31m ",
	"package-lock.json": "\033[31m ",
	"perl":              "\033[33m ",
	"pl":                "\033[33m ",
	"php":               "\033[35m ",
	"py":                "\033[33m ",
	"python":            "\033[33m ",
	"rs":                "\033[31m ",
	"rust":              "\033[31m ",
	"sh":                " ",
	"shell":             " ",
	"sql":               "\033[34m ",
	"swift":             "\033[31m ",
	"toml":              " ",
	"ts":                "\033[34m ",
	"typescript":        "\033[34m ",
	"txt":               " ",
	"text":              " ",
	"xul":               " ",
	"xml":               " ",
	"xhtml":             "\033[33m ",
	"yml":               " ",
	"yaml":              " ",
}

func GetDirIcon(s string) string {
	if i, ok := dirIcon[s]; ok {
		return i
	} else {
		return defaultDir
	}
}

const defaultDir = "\033[96m "

var dirIcon map[string]string = map[string]string{
	"Pictures":     "\033[1;95m󰉏 ",
	"Downloads":    "\033[1;92m󰉍 ",
	".git":         "\033[1;31m ",
	".ssh":         "\033[1;96m󰢬 ",
	"Music":        "\033[1;95m󱍙 ",
	"Desktop":      "\033[1;94m ",
	".vscode":      "\033[1;36m󰨞 ",
	".config":      "\033[1;96m ",
	"config":       "\033[1;96m ",
	"configs":      "\033[1;96m ",
	"bin":          "\033[1;96m ",
	"github":       "\033[1;36m ",
	".github":      "\033[1;36m ",
	"Videos":       "\033[1;95m󰃽 ",
	".cache":       "\033[1;96m󰴌 ",
	".chant":       "\033[1;33m󰝱 ",
	".bard":        "\033[1;33m󰝱 ",
	"node_modules": "\033[1;32m ",
}
