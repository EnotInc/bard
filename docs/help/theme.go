package help

const Theme = `# Theme
Theme contains 3 parts:
	1. General - used across the Bard
	2. Markdown - used in markdown render
	3. Code - used in code render

In each of theme you can provide ASCII escape sequance with color and text type.
For example you can make bold red numbers with ` + "`\\033[1;31m`" + `

To change theme, you can run ` + "`:theme <name>.json`" + ` in bard, or change it in config
You can also get current theme name by running ` + "`:theme`" + ` command
`
