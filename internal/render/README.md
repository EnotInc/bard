# render

So this is a basic renderer. I use a simple `lexer.go` to tokenize text and then style it in `render.go`
Render will have a fiew other renders inside, such as markdonw render, which is located at `markdown` folder

The renderer can process one line at a time, and it keeps current render mode in Renderer.mode as enums.Render type

`cache.go` is used to save rendered lines, so I don't need to render them again if nothing has changed
Buffer saves all lines, except the one with the cursor on it, and only in render. Empty lines and the bottom information line are not saved
