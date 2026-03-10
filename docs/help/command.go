package help

var Command = `# Command Mode

To get back to ***Normal*** mode press <ESC>

## List of commands:
 - q - Exit current buffer. Changes will now be saved
 - qa - Exit all buffers. Changes will now be saved
 - w <file> - saves file or creates a new one if <file> provided
 - wq/x - save buffer and quit
 - h[elp] <topic> - open help file with topic as read only tab
 - rln - switch relative line numeration
 - showmd - show or hide markdown symbols
 - render/rnd - switch markdown render
 - tabnames/tn - show or hide tabs names
 - gt - move to next tab
 - gT - move to previous tab
 - newtab/nt <arg> - creates a new tab. Creates a file if <arg> was provided
`
