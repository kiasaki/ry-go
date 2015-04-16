# ry

_A terminal based, vi like, text editor recreation experiment_

Built with golang, using nsf/termbox-go

## Features

- [x] Opening file from command line
- [x] Display buffer contents
- [x] Windows have a footer with cursor position, modes and buffer name
- [x] Normal mode: by default, receive key events
- [ ] Commands: have a map of command name to functions
- [ ] Commands: support ":" + command + <CR> execution
- [ ] Commands: basic movement (h, j, k, l, 0, ^, $, ^D, ^U)
- [ ] Commands: exit and open buffer
- [ ] Insert mode: add esc support
- [ ] Insert mode: normal keys write to buffer line
- [ ] Insert mode: add enter, tab, del, backspace support
- [ ] Commands & Windows: Support for multiple windows and layouting
- [ ] Commands: save buffer, force quit buffer
- [ ] Commands: shell execution

