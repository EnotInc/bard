# editor

So here is a lot to go through

## `buffer.go`
This file is the main storage for your data

There you can find:
- `line` struct, which is simply a list of runes
- `copied` - which contains data of copied lines
- `cursor` - position of the cursor
- `buffer` - list of lines, copied lines, and 2 cursors. One is the 'real' cursor, and the second is used for visual selection

This file also has all functions to work with data

## `editor.go`
So this is the `editor`. Here is where the terminal is set to RAW mode, and here is the main program loop - func `Run()`

## mode files
For each vim mode that **Bard** supports I use a separate file. Each has a main `caseMode` function, which uses a switch-case construction to decide what to do

## `ui.go`
This file makes one long string, which appears as the editor
`UI` structure contains information about visual cursor, x and y scroll, and render

`Draw` function goes line by line in the buffer and draws it

## `ui_movement.go`
Here is where visual cursor position is changed

