# TODO:
- [x] delete line when backspace was hit and there is no text in line
- [ ] saving files
- [ ] open files
- [ ] add aditinal data of cursor pos
- [ ] add some vim basic command
- [ ] figure out how to work with tabs
- [ ] figure out how to wrap or show long lines of text
- [ ] show and move text up and down
- [ ] render md
- [ ] add file manager in editor(could be side bar or just type over all screen)
- [x] wrap curror at the end of the lines

# BUGS:
- [ ] cursor didn't jump at the bottom of the screen when ':' is pressed
- [ ] cursor style didn't change right away when mode is changed
- [x] fix 'O'

# Refactor:
- [x] rename Program into the Editor
- [?] split into files

# Vim commands:
## Insert:
- [x] 'esc' to normal
- [ ] 'ctrl-c' to normal

## Normal:
- [x] 'i' to insert before
- [x] 'a' to insert after
- [ ] 'I' to insert in the start of the line
- [ ] 'A' to insert in the end of the line
- [x] 'o' to insert a new line bellow
- [x] 'O' to insert a new line above
- [x] ':' to command
- [ ] 'v' to visual
- [ ] 'V' to LineVisual
- [ ] 'x' to delete and copy char under cursor

## Command:
- [x] 'q' to quit
- [ ] 'w' to write
- [ ] 'wq' to write
- [ ] 'x' to save and quit

## Visual:
- [ ] 'x' to delete and copy chars under selectet region
- [ ] 'd' to delete and copy chars under selectet region

## LineVisual:
- [ ] 'x' to delete and copy chars under selectet region
- [ ] 'd' to delete and copy chars under selectet region

