package help

const Theme = `# Theme
Theme contains 3 parts:
	1. General - used across the Bard
	2. Markdown - used in markdown render
	3. Code - used in code render

In each theme you can provide a hex color code without transparency (it will be ignored).
For exmaple code #FF0000 gives you a red color

To change theme, you can run ` + "`:theme <name>.json`" + ` in bard, or change it in config
You can also get current theme name by running ` + "`:theme`" + ` command

To hot reload themes type ` + "`:theme reload`" + `
`
