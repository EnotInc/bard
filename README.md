![Bard Logo](docs/assets/BardLogo256.png)

# **Bard**
## What is Bard
Bard is a little TUI text editor, written in Go with ***0*** third-party dependencies.
I wrote Bard because I wanted to have beautiful markdown rendering, similar to Obsidian, but in the terminal with Vim's motions efficiency.

https://github.com/user-attachments/assets/6f878b12-cee0-45bd-8992-20d79b21f27b

To get general help, you can run Bard with the `-h`/`--help` flag.

## Features
### Explorer
You can toggle focus between the editor and file explorer tiles by pressing `<ctrl+o>`, or turn it on and off with `<ctrl+j>`.

https://github.com/user-attachments/assets/067044ad-fb06-44ff-bab0-aad12069e993


### Editor with vim keys
Supported modes:
- ***Normal***
- ***Command***
- ***Visual***
- ***Visual-Line***

More info about every mode you can find by running `:h[elp] <mode>` command in Bard.

### Space
Space is your own place where you can store your files. You can open it from anywhere by running Bard with the `-s`/`--space` flag.
All files are stored at `~/.bard/space`.

https://github.com/user-attachments/assets/347a0656-0f31-4208-8776-68df3133f8f6

### Other
***Undo/Redo and Command history***

https://github.com/user-attachments/assets/f8df232c-0007-4ad3-8dfc-667c84f0fe7c

---
***Code block render (both for markdown and code files)***

https://github.com/user-attachments/assets/47c3ea71-d33a-48c0-ad12-df87553bdba6

---
***Continuous lists***
When creating a new line (with `o` or `<enter>`) on a list line, Bard will add a new line with the current list type.

https://github.com/user-attachments/assets/bca740f6-9a5d-4f16-82a4-01a9132f95ac

## Settings
At first run, Bard will create the directory `~/.bard/`:
```
.log			// all logs are stored here
config.json		// config file
space/			// your personal Space!
themes/			// theme directory
	bard.json	// default Bard theme
	foo.json	// you can create your own theme
```

You can read more about this in Bard by running the following commands:
- `:h[elp] config`
- `:h[elp] space`
- `:h[elp] theme`

## Installation
Requirements:
- [NerdFont](https://nerdfonts.com). If your terminal doesn't support Nerd Font, you can turn icons off with `:si`/`:showicon` commands in Bard.
- Go 1.25.5.

Bard is not available in any package managers yet, so to use it you have to build it:
```bash
git clone https://github.com/EnotInc/bard.git
cd bard/cmd/bard
go install
```

`go install` will build the project and add `bard` to the `$PATH`, so you don't have to worry about it.
