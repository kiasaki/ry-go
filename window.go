package main

type Window struct {
	Buffer  *Buffer
	CursorX int
	CursorY int
	Editor  *Editor
}

func NewWindow(e *Editor, b *Buffer) *Window {
	return &Window{
		Buffer:  b,
		CursorX: 0,
		CursorY: 0,
		Editor:  e,
	}
}
