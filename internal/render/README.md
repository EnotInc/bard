# render

The renderer can process one line at a time, and it keeps current render mode in Renderer.mode as enums.Render type
Each render mode calls his own render (Code or Markdown for now), and in each directories for those render modes I have a simple lexer, list of tokens and `render.go` file, which turns data from tokens into beautiful text with some ascii escate sequances

`cache.go` is used to save rendered lines, so I don't need to render them again if nothing has changed
Buffer saves all lines, except the one with the cursor on it, and only in render. Empty lines and the bottom information line are not saved
