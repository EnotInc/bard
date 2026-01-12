# TODO:
- [x] delete line when backspace was hit and there is no text in line
- [x] saving files
- [x] open files
- [x] add aditinal data of cursor pos in the end of the line
- [ ] add some vim basic command
- [?] figure out how to work with tabs
- [ ] figure out how to wrap or show long lines of text
- [x] add text scrolling up and down
- [ ] render md
- [ ] add file manager in editor(could be side bar or just type over all screen)
- [x] wrap curror at the end of the lines
- [ ] recalculate window if font was changed

# BUGS:
- [x] 'x' in empty file crashing the program
- [x] line 1000 crashing the editor, "Slice bound out of range [:-1]". Breaking somewhere at buildNumber func (59 line)

# Vim commands:
- [ ] build commands modifyers (like 12j to move 12 line down)

## Insert:
- [x] 'esc' to normal
- [ ] 'ctrl-c' to normal

## Normal:
- [x] 'i' to insert before
- [x] 'a' to insert after
- [?] 'I' to insert in the start of the line
- [x] 'A' to insert in the end of the line
- [x] 'o' to insert a new line bellow
- [x] 'O' to insert a new line above
- [x] ':' to command
- [ ] 'v' to visual
- [ ] 'V' to LineVisual
- [?] 'x' to delete and copy char under cursor
- [x] 's' to delete char under the cursor and enter insert mode
- [ ] 'r' to replace key under the cursor with a new one
- [ ] '.' to repeat last command
- [ ] 'G' to go the end of the file
- [ ] 'gg' to go the beggining of the file


## Command:
- [x] 'q' to quit
- [x] 'w' to write
- [x] 'wq' to write
- [x] 'x' to save and quit
- [x] 'rln' to change line numeration

## Visual:
- [ ] 'x' to delete and copy chars under selectet region
- [ ] 'd' to delete and copy chars under selectet region

## LineVisual:
- [ ] 'x' to delete and copy chars under selectet region
- [ ] 'd' to delete and copy chars under selectet region

# Render
- [x] color cursor and cur line number as yellow
- [?] hide md symbols when cursor is not on the line
- [?] show md symbold when cursor is on the line

