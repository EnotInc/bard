package help

var Config = `# **Config**

Config is located at ` + "`~/.bard/config.json`" + ` file. You can edit it manually, or with commands inside bard

Default config looks like this:
` + "```json" + `
{
` + "\t" + `theme_name: "bard.json" ` + "     " + `# Used to set theme. This field is required '.josn' at the end. All themes is stored at '~/.bard/themes' directory.
` + "\t" + `tab_stop: 4 ` + "                 " + `# Can't be less thant 1. Used to set max tab width.
` + "\t" + `resize_time_duration: 200 ` + "   " + `# Time in milliceconds. Used to set timer to handle terminal resize. Can't be less than 200 and greater than 1000 milliseconds
` + "\t" + `relative_line_number: false ` + " " + `# Used to turn on or off relative line numeration.
` + "\t" + `show_markdown_symbols: false ` + "" + `# If true - always show markdonw symbols (like starts or underlines). By default they are hidden.
` + "\t" + `show_empty_line_symbol: true ` + "" + `# Turns on and off empty line symbol '~'. In read only buffers is off by default.
` + "\t" + `enable_render: true ` + "         " + `# Turns render on and off.
` + "\t" + `show_tab_names: true ` + "        " + `# Used to toggle tabs information in status bar. true: [1|filename], false: [1].
` + "\t" + `keep_tabs: true ` + "             " + `# True by default. If false - replaces inserted tabs with spaces. Doesn't replace tabs in opened file.
` + "\t" + `show_icons: true` + "             " + `# Used in render do decide should it draw nerdfont icons or not.
` + "\t" + `show_borders: true` + "           " + `# Turns on and off tile borders.
` + "\t" + `show_dot_files: true` + "         " + `# Used in explorer to show or hide hidden (dot) files.
}
` + "```" + `
`
