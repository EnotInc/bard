package help

var About = `# **Bard help**

So Bard is vim-like TUI text editor, so if you know some vim motions you'll be fine

Anyways, there are some modes available:
- ***NORMAL*** - main mode where you can move around and do some stuff
- ***INSERT*** - this is where you insert text
- ***COMMAND*** - mode for command execution, like *:wq*
- ***VISUAL*** / *VISUAL-LINE* - mode for text selection

When you run bart you start in *NORMAL* mode, and in any other mode you can press *ESC* to get back to here
For more information run ` + "`" + `:h[elp] <topic>` + "`" + `, for example ` + "`" + `:h command` + "`" + ` to get list of available command in ***Command*** mode

Anyways, here is some main vim keys.
You can also find more information about some specific mode by running *:help <mode>* (not yet tho)

## **Main keys**
### **Movement**
just like in vim:
- ` + "`" + `h` + "`" + ` - move left
- ` + "`" + `j` + "`" + ` - move down
- ` + "`" + `k` + "`" + ` - move up
- ` + "`" + `l` + "`" + ` - move right

You can also type something like *12k* to move 12 lines up

Then you have:
- ` + "`" + `G` + "`" + `- Move to the last line of the file
- ` + "`" + `gg` + "`" + ` - Move to the first line of the file

### **Commands**:
- ` + "`" + `:q` + "`" + ` - quit
- ` + "`" + `:w` + "`" + ` - save (write)
- ` + "`" + `:x` + "`" + ` / ` + "`" + `:wq` + "`" + ` - quit and save (write)

### **How do you type?**
You can type in *INSERT* mode. Here is how you can enter it:
- ` + "`" + `i` + "`" + ` - set mode to ***INSERT*** before the cursor
- ` + "`" + `a` + "`" + ` - set mode to ***INSERT*** after the cursor
- ` + "`" + `I` + "`" + ` - set mode to ***INSERT*** at the first char of the line
- ` + "`" + `A` + "`" + ` - set mode to ***INSERT*** at the end of the line

And here is some helpful vim motions:
- ` + "`" + `o` + "`" + ` - make a new line below and set mode to ***INSERT***
- ` + "`" + `O` + "`" + ` - make a new line above and set mode to ***INSERT***
`
