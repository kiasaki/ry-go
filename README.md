# ry

_a terminal based, experimental text editor, with **vi**'s modal editing and **emacs** extensibility/lispyness_

## Introduction

Built with golang, using nsf/termbox-go. **ry** is the best of both worlds all
with out the weight of 30 years of life. Yes, it's new code, yes it's not as proven
and yes it doesn't support your amiga, but, given the time and the right people
it can be part of the future. Text editor world and tooling hasn't moved much,
but has _bitrotten_ a lot, in the past 30 years and **ry** aims at changing that
(I am looking at make, emacs, vim).

**ry** combines the _modal editing_ and _composability_ of commands from vi that
makes our fingers happy and our productivity go up and the _extensibility_ and
_lispyness_ of emacs.

## Getting started

```
$ go get github.com/kiasaki/ry/cmd/ry
$ ry
```

## Todo list

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

