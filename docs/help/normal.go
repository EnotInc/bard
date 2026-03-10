package help

var Noraml = `# Normal Mode

## Basic vim motion:
 - h, j, k, l - move left, down, up, right
 - you can also type something like 12k to move 12 lines up
 - x - delete char under the cursor
 - G - move to the end of file
 - gg - move to the end of file
 - r - replace 1 char
 - f/F - find char after/before and move to it
 - t/T - find char after/before amd move in front of it
 - w/e/b - work in progress...

## Ways to change modes:
 - i/a - set mode to ***Insert*** before or after the cursor
 - I/A - set mode to ***Insert*** at the srart of at the end of line
 - s - delede char under the cursor and set mode to ***Insert***
 - S - clear line and set mode to ***Insert***
 - o/O - create new line below or above current and set mode to ***Insert**
 - : - set mode to ***Command***
 - v - set mode to ***Visual***
 - V - set mode to ***Visual-line***
 - R - set mode to ***Replace***

## Paste:
 - p/P - paste copied text after/before the cursor
`
