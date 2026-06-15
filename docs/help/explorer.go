package help

var Explorer = `# **Exploror**

To open file explorer press ` + "`<ctrl+j>`" + `. This button combination also hides it if it's already opened.
You can change focused tile without hiding file explorer with ` + "`<ctrl+o>`" + `.

## List of keys:
- ` + "`<esc>`" + ` - switch focus to the editor tile
- ` + "`k`" + ` - move cursor up
- ` + "`j`" + ` - move cursor down
- ` + "`<enter>`" + ` - open file or change directory
- ` + "`d`" + ` - delete file. You must confirm this by pressings ` + "`<enter>`" + `. This will run ` + "`:del <file>`" + ` command
- ` + "`o`" + ` - create a new file or dir. Press ` + "`/`" + ` to change type, ` + "`<enter>`" + ` to create or ` + "`<esc>`" + ` to cancel
- ` + "`g/G`" + ` - move to the top/bottom of the files
`
