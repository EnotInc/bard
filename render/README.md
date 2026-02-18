# render

So this is a basic renderer. I use a simple `lexer.go` to tokenize text and then style it in `render.go`

The renderer can process one line at a time, and right now it can't work with multi-line markdown stuff (like code blocks)
I'll work on this in the near future

`buffer.go` is used to save rendered lines, so I don't need to render them again if nothing has changed
Buffer saves all lines, except the one with the cursor on it, and only in render. Empty lines and the bottom information line are not saved
