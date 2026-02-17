# **Bard**
## **What is Bard**
Bard is a little TUI text editor for markdown files, with vim-like keybindings
I wrote Bard because I wanted to have beautiful markdown rendering, similar to Obsidian, but in the terminal with Vim's motions efficiency

## **Usage**
To run Bard, just type `bard` in the terminal. You can provide a file name to open it or create a new one
To navigate through text, you must use Vim motions, and if you are not familiar with any of this, it's a good time to learn so you can say "*I use vim, btw*" (and Bard, ofc)

## **Installation**
It is not available in any package managers yet, so to install Bard you can do this:
```bash
git clone https://github.com/EnotInc/bard.git
cd bard/cmd/bard
go install
```

The `go install` command will build the project and add `bard` to the `$PATH`
And, yeah, you need `go` to build this project
