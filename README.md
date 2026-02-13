# **Bard**
## **What is bard**
Bard is a little TUI text editor for markdown files, with a vim-like keybindings
I wrote Bard because I wanted to have beautiful markdown rendering, similar to Obsidian, but in terminal with Vim's motions efficiency

## **Usage**
To run bard, just type `bard` in terminal. You can provide file name to open it, or create a new one.
To navigate through text, you must use Vim motions, and if you not familiar with any of this, it's a good time to lear so you can say "*I use vim, btw*" (and Bard, ofc)

## **Instalation**
It is not available in any package managers yet, so to install Bard you can do this:
```bash
git clone https://github.com/EnotInc/bard.git
cd bard/cmd/bard
go install
```

The `go install` command will build project and add `Bard` to the `$PATH`
And, yeah, you need `go` to build this project

## **Supported modes**
- `Normal`
- `Insert`
- `Command`

